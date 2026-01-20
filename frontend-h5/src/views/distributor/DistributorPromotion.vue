<template>
   <div class="distributor-promotion">
     <van-nav-bar title="推广工具" left-arrow @click-left="$router.back()" />
 
     <van-tabs v-model:active="activeTab" sticky>
       <!-- 活动专属海报 -->
       <van-tab title="活动海报" name="campaign">
         <div class="promotion-content">
           <!-- 活动选择 -->
           <van-cell-group inset>
             <van-field
               v-model="selectedCampaignId"
               name="campaign"
               label="选择活动"
               placeholder="请选择要推广的活动"
               is-link
               readonly
               @click="showCampaignPicker = true"
             >
               <template #input>
                 {{ selectedCampaignName || '请选择活动' }}
               </template>
             </van-field>
           </van-cell-group>
 
           <!-- 海报生成按钮 -->
           <div class="poster-actions">
             <van-button
               block
               type="primary"
               icon="photo-o"
               @click="generateCampaignPoster"
               :loading="generatingPoster"
             >
               生成活动专属海报
             </van-button>
           </div>
 
           <!-- 海报预览 -->
           <div class="poster-preview" v-if="campaignPoster">
             <h3>海报预览</h3>
             <div class="poster-image-container">
               <img :src="campaignPoster.templateUrl" alt="活动海报" v-if="campaignPoster.templateUrl" />
               <van-loading v-else-if="generatingPoster" size="large" />
               <div class="poster-placeholder" v-else>
                 <van-empty description="请选择活动并生成海报" />
               </div>
             </div>
 
             <!-- 海报信息 -->
             <van-cell-group inset class="poster-info">
               <van-cell title="活动名称" :value="campaignPosterData.campaignName" />
               <van-cell title="活动描述" :value="campaignPosterData.campaignDesc" />
               <van-cell title="分销商" :value="campaignPosterData.distributorName" />
             </van-cell-group>
 
             <!-- 操作按钮 -->
             <div class="action-buttons">
               <van-button
                 block
                 type="primary"
                 icon="down"
                 @click="downloadPoster(campaignPoster.templateUrl, '活动海报')"
               >
                 下载海报
               </van-button>
               <van-button
                 block
                 icon="share-o"
                 @click="sharePoster(campaignPoster.templateUrl)"
               >
                 分享海报
               </van-button>
             </div>
           </div>
         </div>
       </van-tab>
 
       <!-- 通用分销商海报 -->
       <van-tab title="分销商海报" name="distributor">
         <div class="promotion-content">
           <!-- 海报生成按钮 -->
           <div class="poster-actions">
             <van-button
               block
               type="primary"
               icon="photo-o"
               @click="generateDistributorPoster"
               :loading="generatingPoster"
             >
               生成分销商专属海报
             </van-button>
           </div>
 
           <!-- 海报预览 -->
           <div class="poster-preview" v-if="distributorPoster">
             <h3>海报预览</h3>
             <div class="poster-image-container">
               <img :src="distributorPoster.templateUrl" alt="分销商海报" v-if="distributorPoster.templateUrl" />
               <van-loading v-else-if="generatingPoster" size="large" />
               <div class="poster-placeholder" v-else>
                 <van-empty description="点击生成海报" />
               </div>
             </div>
 
             <!-- 海报信息 -->
             <van-cell-group inset class="poster-info">
               <van-cell title="分销商" :value="distributorPosterData.distributorName" />
               <van-cell title="活动数量" :value="distributorPosterData.campaignCount" />
               <van-cell title="推广说明" value="包含所有可推广活动" />
             </van-cell-group>
 
             <!-- 操作按钮 -->
             <div class="action-buttons">
               <van-button
                 block
                 type="primary"
                 icon="down"
                 @click="downloadPoster(distributorPoster.templateUrl, '分销商海报')"
               >
                 下载海报
               </van-button>
               <van-button
                 block
                 icon="share-o"
                 @click="sharePoster(distributorPoster.templateUrl)"
               >
                 分享海报
               </van-button>
             </div>
           </div>
         </div>
       </van-tab>
 
       <!-- 推广链接 -->
       <van-tab title="推广链接" name="link">
         <div class="promotion-content">
           <!-- 活动选择 -->
           <van-cell-group inset>
             <van-field
               v-model="selectedCampaignId"
               name="campaign"
               label="选择活动"
               placeholder="请选择要推广的活动"
               is-link
               readonly
               @click="showCampaignPicker = true"
             >
               <template #input>
                 {{ selectedCampaignName || '请选择活动' }}
               </template>
             </van-field>
           </van-cell-group>
 
           <!-- 推广链接信息 -->
           <div class="link-info" v-if="generatedLink">
             <van-cell-group inset>
               <van-cell title="推广链接" :value="generatedLink.link" />
               <van-cell title="推广码" :value="generatedLink.linkCode" />
             </van-cell-group>
 
             <!-- 二维码 -->
             <div class="qrcode-section">
               <h3>推广二维码</h3>
               <div class="qrcode-container">
                 <img :src="qrcodeUrl" alt="推广二维码" v-if="qrcodeUrl" />
                 <van-loading v-else-if="loadingQrcode" />
                 <div class="qrcode-placeholder" v-else>
                   <van-button size="small" type="primary" @click="loadQrcode">
                     生成二维码
                   </van-button>
                 </div>
               </div>
               <p class="qrcode-tip">扫描二维码即可访问活动页面</p>
             </div>
 
             <!-- 操作按钮 -->
             <div class="action-buttons">
               <van-button
                 block
                 type="primary"
                 icon="link-o"
                 @click="copyLink"
               >
                 复制推广链接
               </van-button>
               <van-button
                 block
                 icon="down"
                 @click="downloadQrcode"
                 v-if="qrcodeUrl"
               >
                 下载二维码
               </van-button>
             </div>
           </div>
 
           <!-- 我的推广链接 -->
           <div class="my-links">
             <h3>我的推广链接</h3>
             <van-empty v-if="links.length === 0" description="暂无推广链接" />
             <van-cell-group inset v-else>
               <van-cell
                 v-for="link in links"
                 :key="link.linkId"
                 :title="link.campaignName || `活动 ${link.campaignId}`"
                 :value="`点击 ${link.clickCount} 次`"
                 is-link
                 @click="viewLink(link)"
               />
             </van-cell-group>
           </div>
         </div>
       </van-tab>
     </van-tabs>
 
     <!-- 推广说明 -->
     <div class="promotion-guide">
       <h3>推广说明</h3>
       <div class="guide-content">
         <p>1. 选择要推广的活动，生成专属海报或链接</p>
         <p>2. 将海报或二维码分享给朋友，引导他们下单</p>
         <p>3. 用户通过您的推广下单，您即可获得佣金奖励</p>
         <p>4. 最多支持三级分销，您的下级推广也能为您带来收益</p>
       </div>
     </div>
 
     <!-- 活动选择器 -->
     <van-popup v-model:show="showCampaignPicker" position="bottom">
       <van-picker
         :columns="campaignOptions"
         @confirm="onCampaignConfirm"
         @cancel="showCampaignPicker = false"
       />
     </van-popup>
   </div>
 </template>

<script>
 import { ref, computed, onMounted } from 'vue'
 import { useRoute, useRouter } from 'vue-router'
 import { Toast } from 'vant'
 import axios from '@/utils/axios'
 
 export default {
   name: 'DistributorPromotion',
   setup() {
     const route = useRoute()
     const router = useRouter()
     const brandId = ref(parseInt(route.query.brandId) || 0)
     const activeTab = ref('campaign')
     const selectedCampaignId = ref(null)
     const selectedCampaignName = ref('')
     const showCampaignPicker = ref(false)
     const generatedLink = ref(null)
     const qrcodeUrl = ref('')
     const loadingQrcode = ref(false)
     const links = ref([])
     const campaigns = ref([])
     const campaignPoster = ref(null)
     const distributorPoster = ref(null)
     const generatingPoster = ref(false)
     const campaignPosterData = ref({})
     const distributorPosterData = ref({})
 
     const campaignOptions = computed(() => {
       return campaigns.value.map(c => ({
         text: c.name,
         value: c.id
       }))
     })
 
     // 加载活动列表
     const loadCampaigns = async () => {
       try {
         const { data } = await axios.get('/api/v1/campaigns', {
           params: { status: 'active', pageSize: 100 }
         })
         if (data.code === 200) {
           campaigns.value = data.data.campaigns || []
         }
       } catch (error) {
         console.error('获取活动列表失败:', error)
       }
     }
 
     // 加载我的推广链接
     const loadMyLinks = async () => {
       try {
         const { data } = await axios.get('/api/v1/distributor/links')
         if (data.code === 200) {
           links.value = data.data || []
           // 获取活动名称
           for (const link of links.value) {
             const campaign = campaigns.value.find(c => c.id === link.campaignId)
             if (campaign) {
               link.campaignName = campaign.name
             }
           }
         }
       } catch (error) {
         console.error('获取推广链接失败:', error)
       }
     }
 
     // 生成活动专属海报
     const generateCampaignPoster = async () => {
       if (!selectedCampaignId.value) {
         Toast('请先选择活动')
         return
       }

       generatingPoster.value = true
       try {
         const { data } = await axios.post('/api/v1/posters/generate', {
           type: 'campaign',
           campaignId: selectedCampaignId.value
         })
         if (data.code === 200) {
           campaignPoster.value = data.data
           campaignPosterData.value = data.data.posterData || {}
           Toast('海报生成成功')
         }
       } catch (error) {
         Toast(error.response?.data?.message || '海报生成失败')
       } finally {
         generatingPoster.value = false
       }
     }

    // 生成分销商专属海报
     const generateDistributorPoster = async () => {
       generatingPoster.value = true
       try {
         const { data } = await axios.post('/api/v1/posters/generate', {
           type: 'distributor'
         })
         if (data.code === 200) {
           distributorPoster.value = data.data
           distributorPosterData.value = data.data.posterData || {}
           Toast('海报生成成功')
         }
       } catch (error) {
         Toast(error.response?.data?.message || '海报生成失败')
       } finally {
         generatingPoster.value = false
       }
     }
 
     // 下载海报
     const downloadPoster = (url, name) => {
       if (url) {
         const link = document.createElement('a')
         link.href = url
         link.download = `${name}_${Date.now()}.png`
         link.click()
         Toast('正在下载海报...')
       }
     }
 
     // 分享海报
     const sharePoster = async (url) => {
       if (navigator.share && url) {
         try {
           await navigator.share({
             title: '分享我的推广海报',
             url: url
           })
         } catch (error) {
           Toast('分享已取消')
         }
       } else {
         Toast('当前环境不支持分享')
       }
     }
 
     // 选择活动
     const onCampaignConfirm = async ({ selectedOptions }) => {
       const option = selectedOptions[0]
       selectedCampaignId.value = option.value
       selectedCampaignName.value = option.text
       showCampaignPicker.value = false
 
       // 生成推广链接
       await generateLink(option.value)
     }
 
     // 生成推广链接
     const generateLink = async (campaignId) => {
       try {
         const { data } = await axios.post('/api/v1/distributor/link/generate', {
           campaignId: campaignId
         })
         if (data.code === 200) {
           generatedLink.value = data.data
           qrcodeUrl.value = ''
           // 重新加载链接列表
           loadMyLinks()
         }
       } catch (error) {
         Toast(error.response?.data?.message || '生成链接失败')
       }
     }
 
     // 加载二维码
     const loadQrcode = () => {
       loadingQrcode.value = true
       // 使用二维码API
       qrcodeUrl.value = generatedLink.value.qrcodeUrl || `https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(generatedLink.value.link)}`
       loadingQrcode.value = false
     }
 
     // 复制链接
     const copyLink = () => {
       if (generatedLink.value) {
         // 使用 Clipboard API
         if (navigator.clipboard) {
           navigator.clipboard.writeText(generatedLink.value.link)
           Toast('链接已复制')
         } else {
           // 降级方案
           const input = document.createElement('input')
           input.value = generatedLink.value.link
           document.body.appendChild(input)
           input.select()
           document.execCommand('copy')
           document.body.removeChild(input)
           Toast('链接已复制')
         }
       }
     }
 
     // 下载二维码
     const downloadQrcode = () => {
       if (qrcodeUrl.value) {
         const link = document.createElement('a')
         link.href = qrcodeUrl.value
         link.download = `推广二维码_${generatedLink.value.linkCode}.png`
         link.click()
         Toast('正在下载二维码...')
       }
     }
 
     // 查看链接详情
     const viewLink = (link) => {
       selectedCampaignId.value = link.campaignId
       const campaign = campaigns.value.find(c => c.id === link.campaignId)
       selectedCampaignName.value = campaign ? campaign.name : `活动 ${link.campaignId}`
       generatedLink.value = link
       qrcodeUrl.value = ''
     }
 
     onMounted(() => {
       loadCampaigns()
       loadMyLinks()
     })
 
     return {
       brandId,
       activeTab,
       selectedCampaignId,
       selectedCampaignName,
       showCampaignPicker,
       generatedLink,
       qrcodeUrl,
       loadingQrcode,
       links,
       campaigns,
       campaignPoster,
       distributorPoster,
       generatingPoster,
       campaignPosterData,
       distributorPosterData,
       campaignOptions,
       generateCampaignPoster,
       generateDistributorPoster,
       downloadPoster,
       sharePoster,
       onCampaignConfirm,
       loadQrcode,
       copyLink,
       downloadQrcode,
       viewLink
     }
   }
 }
 </script>

<style scoped>
.distributor-promotion {
  min-height: 100vh;
  background-color: #f5f5f5;
}

.promotion-content {
  padding: 16px;
}

.van-cell-group {
  margin-bottom: 16px;
}

.link-info {
  margin-top: 16px;
}

.qrcode-section {
  background: white;
  margin: 16px;
  padding: 16px;
  border-radius: 8px;
  text-align: center;
}

.qrcode-section h3 {
  margin: 0 0 16px;
  font-size: 16px;
  color: #333;
}

.qrcode-container {
  width: 200px;
  height: 200px;
  margin: 0 auto 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  border-radius: 8px;
}

.qrcode-container img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.qrcode-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.qrcode-tip {
  font-size: 12px;
  color: #999;
  margin: 0;
}

.action-buttons {
  padding: 0 16px;
  margin-bottom: 16px;
}

.action-buttons .van-button {
  margin-bottom: 12px;
}

.poster-actions {
  padding: 16px;
}

.poster-preview {
  margin-top: 16px;
  padding: 16px;
  background: white;
  border-radius: 8px;
}

.poster-preview h3 {
  margin: 0 0 16px;
  font-size: 16px;
  color: #333;
}

.poster-image-container {
  width: 100%;
  min-height: 400px;
  margin: 0 auto 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  border-radius: 8px;
  overflow: hidden;
}

.poster-image-container img {
  width: 100%;
  height: auto;
  object-fit: contain;
}

.poster-placeholder {
  width: 100%;
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.poster-info {
  margin-bottom: 16px;
}

.my-links {
  margin-top: 24px;
}

.my-links h3 {
  padding: 0 16px;
  font-size: 16px;
  color: #333;
  margin-bottom: 12px;
}

.promotion-guide {
  margin-top: 24px;
  background: white;
  padding: 16px;
  border-radius: 8px;
}

.promotion-guide h3 {
  margin: 0 0 12px;
  font-size: 16px;
  color: #333;
}

.guide-content p {
  margin: 8px 0;
  font-size: 14px;
  color: #666;
  line-height: 1.6;
}
</style>
