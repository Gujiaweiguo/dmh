<template>
  <div class="distributor-withdrawals">
    <!-- 顶部卡片 -->
    <div class="balance-card">
      <div class="balance-info">
        <div class="balance-label">可提现余额</div>
        <div class="balance-value">¥{{ balance.toFixed(2) }}</div>
        <van-button type="primary" size="small" @click="showWithdrawalDialog = true">
          申请提现
        </van-button>
      </div>
    </div>

    <!-- 提现记录列表 -->
    <div class="withdrawals-list">
      <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
        <van-list
          v-model:loading"
          :finished="finished"
          finished-text="没有更多了"
          @load="onLoad"
        >
          <div
            v-for="withdrawal in withdrawals"
            :key="withdrawal.id"
            class="withdrawal-item"
          >
            <div class="withdrawal-header">
              <div class="withdrawal-amount">
                <span class="amount">¥{{ withdrawal.amount.toFixed(2) }}</span>
                <van-tag :type="getStatusType(withdrawal.status)">
                  {{ getStatusText(withdrawal.status) }}
                </van-tag>
              </div>
              <div class="withdrawal-date">{{ formatDate(withdrawal.createdAt) }}</div>
            </div>
            
            <div class="withdrawal-details">
              <div class="detail-item">
                <span class="label">提现方式：</span>
                <span class="value">{{ getPayTypeText(withdrawal.payType) }}</span>
              </div>
              <div class="detail-item" v-if="withdrawal.payAccount">
                <span class="label">提现账号：</span>
                <span class="value">{{ withdrawal.payAccount }}</span>
              </div>
              <div class="detail-item" v-if="withdrawal.rejectedReason">
                <span class="label">拒绝原因：</span>
                <span class="value rejected">{{ withdrawal.rejectedReason }}</span>
              </div>
            </div>
          </div>
        </van-list>
      </van-pull-refresh>
    </div>

    <!-- 申请提现弹窗 -->
    <van-dialog
      v-model:showWithdrawalDialog
      title="申请提现"
      show-cancel-button
      @confirm="applyWithdrawal"
      @cancel="resetWithdrawalForm"
    >
      <van-form @submit="applyWithdrawal">
        <van-field
          v-model="withdrawalForm.amount"
          name="amount"
          label="提现金额"
          type="number"
          placeholder="请输入提现金额"
          :rules="amountRules"
        />
        <van-field
          v-model="withdrawalForm.payType"
          name="payType"
          label="提现方式"
          readonly
          is-link
          @click="showPayTypePicker = true"
          placeholder="请选择提现方式"
        />
        <van-field
          v-model="withdrawalForm.payAccount"
          name="payAccount"
          :label="getPayAccountLabel()"
          type="text"
          placeholder="请输入提现账号"
          :rules="payAccountRules"
        />
        <van-field
          v-model="withdrawalForm.payRealName"
          name="payRealName"
          :label="getPayRealNameLabel()"
          type="text"
          placeholder="请输入真实姓名"
          :rules="payRealNameRules"
        />
      </van-form>
    </van-dialog>

    <!-- 提现方式选择器 -->
    <van-action-sheet v-model:showPayTypePicker>
      <van-cell
        v-for="payType in payTypes"
        :key="payType.value"
        :title="payType.label"
        is-link
        @click="selectPayType(payType)"
      />
    </van-action-sheet>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Toast, Dialog } from 'vant'
import axios from '@/utils/axios'

export default {
  name: 'DistributorWithdrawals',
  setup() {
    const router = useRouter()
    const balance = ref(0)
    const withdrawals = ref([])
    const loading = ref(false)
    const finished = ref(false)
    const refreshing = ref(false)
    const page = ref(1)
    const pageSize = ref(20)
    const total = ref(0)
    
    const showWithdrawalDialog = ref(false)
    const showPayTypePicker = ref(false)
    
    const withdrawalForm = ref({
      amount: '',
      payType: 'wechat',
      payAccount: '',
      payRealName: ''
    })
    
    const payTypes = [
      { label: '微信', value: 'wechat' },
      { label: '支付宝', value: 'alipay' },
      { label: '银行卡', value: 'bank' }
    ]
    
    // 表单验证规则
    const amountRules = [
      { required: true, message: '请输入提现金额' },
      {
        validator: value => {
          const amount = parseFloat(value)
          if (isNaN(amount) || amount <= 0) {
            return '提现金额必须大于0'
          }
          if (amount > balance.value) {
            return '提现金额不能超过可用余额'
          }
          return true
        }
      }
    ]
    
    const payAccountRules = [
      { required: true, message: '请输入提现账号' }
    ]
    
    const payRealNameRules = [
      { required: true, message: '请输入真实姓名' }
    ]
    
    // 获取提现方式选择器标签
    const getPayAccountLabel = () => {
      const payTypeMap = {
        wechat: '微信号',
        alipay: '支付宝账号',
        bank: '银行卡号'
      }
      return payTypeMap[withdrawalForm.value.payType] || '提现账号'
    }
    
    // 获取真实姓名标签
    const getPayRealNameLabel = () => {
      const payTypeMap = {
        wechat: '微信昵称',
        alipay: '支付宝姓名',
        bank: '银行卡持卡人'
      }
      return payTypeMap[withdrawalForm.value.payType] || '真实姓名'
    }
    
    // 加载提现记录
    const loadWithdrawals = async (isRefresh = false) => {
      if (loading.value || (finished.value && !isRefresh)) {
        return
      }
      
      loading.value = true
      
      try {
        const { data } = await axios.get('/api/v1/withdrawals/my', {
          params: {
            page: page.value,
            pageSize: pageSize.value
          }
        })
        
        if (data.code === 200) {
          if (isRefresh) {
            withdrawals.value = data.data.withdrawals
          } else {
            withdrawals.value.push(...data.data.withdrawals)
          }
          total.value = data.data.total
          balance.value = data.data.balance || 0
          
          finished.value = withdrawals.value.length >= total.value
        }
      } catch (error) {
        console.error('加载提现记录失败:', error)
        Toast.fail('加载失败，请重试')
      } finally {
        loading.value = false
        refreshing.value = false
      }
    }
    
    // 刷新
    const onRefresh = async () => {
      page.value = 1
      finished.value = false
      await loadWithdrawals(true)
    }
    
    // 加载更多
    const onLoad = async () => {
      if (finished.value) {
        return
      }
      page.value++
      await loadWithdrawals()
    }
    
    // 申请提现
    const applyWithdrawal = async () => {
      try {
        const { data } = await axios.post('/api/v1/withdrawals/apply', {
          amount: parseFloat(withdrawalForm.value.amount),
          payType: withdrawalForm.value.payType,
          payAccount: withdrawalForm.value.payAccount,
          payRealName: withdrawalForm.value.payRealName
        })
        
        if (data.code === 200) {
          Toast.success('提现申请已提交')
          showWithdrawalDialog.value = false
          resetWithdrawalForm()
          // 重新加载提现记录
          page.value = 1
          finished.value = false
          withdrawals.value = []
          await loadWithdrawals()
          // 更新余额
          balance.value -= parseFloat(withdrawalForm.value.amount)
        }
      } catch (error) {
        console.error('申请提现失败:', error)
        Toast.fail('申请失败，请重试')
      }
    }
    
    // 选择提现方式
    const selectPayType = (payType) => {
      withdrawalForm.value.payType = payType.value
      showPayTypePicker.value = false
    }
    
    // 重置表单
    const resetWithdrawalForm = () => {
      withdrawalForm.value = {
        amount: '',
        payType: 'wechat',
        payAccount: '',
        payRealName: ''
      }
    }
    
    // 获取状态类型
    const getStatusType = (status) => {
      const types = {
        pending: 'warning',
        approved: 'primary',
        rejected: 'danger',
        processing: 'success',
        completed: 'success',
        failed: 'danger'
      }
      return types[status] || 'default'
    }
    
    // 获取状态文本
    const getStatusText = (status) => {
      const texts = {
        pending: '待审核',
        approved: '已批准',
        rejected: '已拒绝',
        processing: '处理中',
        completed: '已完成',
        failed: '失败'
      }
      return texts[status] || status
    }
    
    // 获取提现方式文本
    const getPayTypeText = (payType) => {
      const texts = {
        wechat: '微信',
        alipay: '支付宝',
        bank: '银行卡'
      }
      return texts[payType] || payType
    }
    
    // 格式化日期
    const formatDate = (date) => {
      if (!date) return ''
      const d = new Date(date)
      return `${d.getFullYear()}-${(d.getMonth() + 1).toString().padStart(2, '0')}-${d.getDate().toString().padStart(2, '0')} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}`
    }
    
    onMounted(() => {
      loadWithdrawals()
    })
    
    return {
      balance,
      withdrawals,
      loading,
      finished,
      refreshing,
      page,
      pageSize,
      total,
      showWithdrawalDialog,
      showPayTypePicker,
      withdrawalForm,
      payTypes,
      amountRules,
      payAccountRules,
      payRealNameRules,
      getPayAccountLabel,
      getPayRealNameLabel,
      selectPayType,
      resetWithdrawalForm,
      applyWithdrawal,
      onRefresh,
      onLoad,
      getStatusType,
      getStatusText,
      getPayTypeText,
      formatDate
    }
  }
}
</script>

<style scoped>
.distributor-withdrawals {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 50px;
}

.balance-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 30px 20px;
  margin-bottom: 16px;
}

.balance-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.balance-label {
  font-size: 16px;
  opacity: 0.9;
}

.balance-value {
  font-size: 32px;
  font-weight: bold;
}

.withdrawals-list {
  background: white;
}

.withdrawal-item {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.withdrawal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.withdrawal-amount {
  display: flex;
  align-items: center;
  gap: 8px;
}

.amount {
  font-size: 20px;
  font-weight: bold;
  color: #333;
}

.withdrawal-date {
  font-size: 12px;
  color: #999;
}

.withdrawal-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-item {
  display: flex;
  font-size: 14px;
}

.detail-item .label {
  color: #999;
  width: 80px;
  flex-shrink: 0;
}

.detail-item .value {
  color: #333;
  flex: 1;
}

.detail-item .value.rejected {
  color: #f56c6c;
}
</style>
