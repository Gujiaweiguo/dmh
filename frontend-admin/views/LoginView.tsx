import { defineComponent, h, ref, reactive, computed } from 'vue';
import * as LucideIcons from 'lucide-vue-next';
import { authApi } from '../services/authApi';
import { LoginRequest, RegisterRequest } from '../types';

// 公开接口类型定义 (Vue component instance 自动解包 Ref/Reactive)
export interface LoginViewInstance {
  isLogin: boolean;
  loading: boolean;
  error: string;
  loginForm: LoginRequest;
  registerForm: RegisterRequest;
  loginValid: boolean;
  registerValid: boolean;
  handleLogin: () => Promise<void>;
  handleRegister: () => Promise<void>;
  toggleMode: () => void;
}

export const LoginView = defineComponent({
  emits: ['login-success'],
  setup(props, { emit }) {
    const isLogin = ref(true);
    const loading = ref(false);
    const error = ref('');
    
    // 登录表单
    const loginForm = reactive<LoginRequest>({
      username: '',
      password: ''
    });
    
    // 注册表单
    const registerForm = reactive<RegisterRequest>({
      username: '',
      password: '',
      phone: '',
      email: '',
      realName: ''
    });
    
    // 表单验证
    const loginValid = computed(() => {
      return loginForm.username.length > 0 && loginForm.password.length >= 6;
    });
    
    const registerValid = computed(() => {
      return registerForm.username.length > 0 && 
             registerForm.password.length >= 6 && 
             registerForm.phone.length > 0 &&
             /^1[3-9]\d{9}$/.test(registerForm.phone);
    });
    
    // 处理登录
    const handleLogin = async () => {
      if (!loginValid.value) {
        error.value = '请填写完整的登录信息';
        return;
      }
      
      loading.value = true;
      error.value = '';
      
      try {
        const result = await authApi.login(loginForm);
        emit('login-success', result);
      } catch (err: any) {
        error.value = err.message || '登录失败，请检查用户名和密码';
      } finally {
        loading.value = false;
      }
    };
    
    // 处理注册
    const handleRegister = async () => {
      if (!registerValid.value) {
        error.value = '请填写完整且正确的注册信息';
        return;
      }
      
      loading.value = true;
      error.value = '';
      
      try {
        const result = await authApi.register(registerForm);
        emit('login-success', result);
      } catch (err: any) {
        error.value = err.message || '注册失败，请重试';
      } finally {
        loading.value = false;
      }
    };
    
    // 切换登录/注册模式
    const toggleMode = () => {
      isLogin.value = !isLogin.value;
      error.value = '';
    };

    return () => h('div', { class: 'min-h-screen bg-gradient-to-br from-indigo-50 via-white to-purple-50 flex items-center justify-center p-4' }, [
      h('div', { class: 'w-full max-w-md' }, [
        // Logo和标题
        h('div', { class: 'text-center mb-8' }, [
          h('div', { class: 'w-20 h-20 bg-gradient-to-br from-indigo-600 to-purple-600 rounded-3xl flex items-center justify-center mx-auto mb-6 shadow-lg' }, [
            h(LucideIcons.Zap, { size: 32, class: 'text-white' })
          ]),
          h('h1', { class: 'text-3xl font-black text-slate-900 mb-2' }, 'DMH 数字营销中台'),
          h('p', { class: 'text-slate-500' }, isLogin.value ? '登录您的管理账户' : '创建新的管理账户')
        ]),
        
        // 登录/注册表单
        h('div', { class: 'bg-white rounded-3xl shadow-xl border border-slate-100 p-8' }, [
          // 模式切换标签
          h('div', { class: 'flex bg-slate-100 rounded-2xl p-1 mb-8' }, [
            h('button', {
              onClick: () => { isLogin.value = true; error.value = ''; },
              class: `flex-1 py-3 px-4 rounded-xl font-bold text-sm transition-all ${
                isLogin.value 
                  ? 'bg-white text-slate-900 shadow-sm' 
                  : 'text-slate-500 hover:text-slate-700'
              }`
            }, '登录'),
            h('button', {
              onClick: () => { isLogin.value = false; error.value = ''; },
              class: `flex-1 py-3 px-4 rounded-xl font-bold text-sm transition-all ${
                !isLogin.value 
                  ? 'bg-white text-slate-900 shadow-sm' 
                  : 'text-slate-500 hover:text-slate-700'
              }`
            }, '注册')
          ]),
          
          // 错误提示
          error.value && h('div', { class: 'mb-6 p-4 bg-red-50 border border-red-200 rounded-2xl flex items-start gap-3' }, [
            h(LucideIcons.AlertCircle, { size: 20, class: 'text-red-600 flex-shrink-0 mt-0.5' }),
            h('div', { class: 'text-sm text-red-800' }, error.value)
          ]),
          
          // 登录表单
          isLogin.value && h('form', { 
            onSubmit: (e: Event) => { e.preventDefault(); handleLogin(); },
            class: 'space-y-6' 
          }, [
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '用户名'),
              h('div', { class: 'relative' }, [
                h('div', { class: 'absolute left-4 top-1/2 -translate-y-1/2 text-slate-400' }, 
                  h(LucideIcons.User, { size: 20 })
                ),
                h('input', {
                  type: 'text',
                  value: loginForm.username,
                  onInput: (e: any) => { loginForm.username = e.target.value },
                  placeholder: '请输入用户名',
                  class: 'w-full pl-12 pr-4 py-4 rounded-2xl border border-slate-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all',
                  required: true
                })
              ])
            ]),
            
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '密码'),
              h('div', { class: 'relative' }, [
                h('div', { class: 'absolute left-4 top-1/2 -translate-y-1/2 text-slate-400' }, 
                  h(LucideIcons.Lock, { size: 20 })
                ),
                h('input', {
                  type: 'password',
                  value: loginForm.password,
                  onInput: (e: any) => { loginForm.password = e.target.value },
                  placeholder: '请输入密码',
                  class: 'w-full pl-12 pr-4 py-4 rounded-2xl border border-slate-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all',
                  required: true,
                  minlength: 6
                })
              ])
            ]),
            
            h('button', {
              type: 'submit',
              disabled: loading.value || !loginValid.value,
              class: `w-full py-4 rounded-2xl font-bold text-white transition-all ${
                loading.value || !loginValid.value
                  ? 'bg-slate-300 cursor-not-allowed'
                  : 'bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 shadow-lg hover:shadow-xl'
              }`
            }, loading.value ? '登录中...' : '登录')
          ]),
          
          // 注册表单
          !isLogin.value && h('form', { 
            onSubmit: (e: Event) => { e.preventDefault(); handleRegister(); },
            class: 'space-y-6' 
          }, [
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '用户名'),
              h('div', { class: 'relative' }, [
                h('div', { class: 'absolute left-4 top-1/2 -translate-y-1/2 text-slate-400' }, 
                  h(LucideIcons.User, { size: 20 })
                ),
                h('input', {
                  type: 'text',
                  value: registerForm.username,
                  onInput: (e: any) => { registerForm.username = e.target.value },
                  placeholder: '请输入用户名',
                  class: 'w-full pl-12 pr-4 py-4 rounded-2xl border border-slate-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all',
                  required: true
                })
              ])
            ]),
            
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '手机号'),
              h('div', { class: 'relative' }, [
                h('div', { class: 'absolute left-4 top-1/2 -translate-y-1/2 text-slate-400' }, 
                  h(LucideIcons.Phone, { size: 20 })
                ),
                h('input', {
                  type: 'tel',
                  value: registerForm.phone,
                  onInput: (e: any) => { registerForm.phone = e.target.value },
                  placeholder: '请输入手机号',
                  class: 'w-full pl-12 pr-4 py-4 rounded-2xl border border-slate-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all',
                  required: true,
                  pattern: '^1[3-9]\\d{9}$'
                })
              ])
            ]),
            
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '真实姓名'),
              h('div', { class: 'relative' }, [
                h('div', { class: 'absolute left-4 top-1/2 -translate-y-1/2 text-slate-400' }, 
                  h(LucideIcons.UserCheck, { size: 20 })
                ),
                h('input', {
                  type: 'text',
                  value: registerForm.realName,
                  onInput: (e: any) => { registerForm.realName = e.target.value },
                  placeholder: '请输入真实姓名（可选）',
                  class: 'w-full pl-12 pr-4 py-4 rounded-2xl border border-slate-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all'
                })
              ])
            ]),
            
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '邮箱'),
              h('div', { class: 'relative' }, [
                h('div', { class: 'absolute left-4 top-1/2 -translate-y-1/2 text-slate-400' }, 
                  h(LucideIcons.Mail, { size: 20 })
                ),
                h('input', {
                  type: 'email',
                  value: registerForm.email,
                  onInput: (e: any) => { registerForm.email = e.target.value },
                  placeholder: '请输入邮箱（可选）',
                  class: 'w-full pl-12 pr-4 py-4 rounded-2xl border border-slate-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all'
                })
              ])
            ]),
            
            h('div', [
              h('label', { class: 'block text-sm font-bold text-slate-700 mb-2' }, '密码'),
              h('div', { class: 'relative' }, [
                h('div', { class: 'absolute left-4 top-1/2 -translate-y-1/2 text-slate-400' }, 
                  h(LucideIcons.Lock, { size: 20 })
                ),
                h('input', {
                  type: 'password',
                  value: registerForm.password,
                  onInput: (e: any) => { registerForm.password = e.target.value },
                  placeholder: '请输入密码（至少6位）',
                  class: 'w-full pl-12 pr-4 py-4 rounded-2xl border border-slate-200 outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-200 transition-all',
                  required: true,
                  minlength: 6
                })
              ])
            ]),
            
            h('button', {
              type: 'submit',
              disabled: loading.value || !registerValid.value,
              class: `w-full py-4 rounded-2xl font-bold text-white transition-all ${
                loading.value || !registerValid.value
                  ? 'bg-slate-300 cursor-not-allowed'
                  : 'bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-500 hover:to-purple-500 shadow-lg hover:shadow-xl'
              }`
            }, loading.value ? '注册中...' : '注册账户')
          ]),
          
          // 底部提示
          h('div', { class: 'mt-8 text-center text-sm text-slate-500' }, [
            // 测试账号快速填充
            !isLogin.value ? null : h('div', { class: 'mb-4 p-4 bg-amber-50 border border-amber-200 rounded-2xl' }, [
              h('p', { class: 'text-amber-800 font-bold mb-2' }, '⚠️ 测试账号'),
              h('div', { class: 'text-amber-700 text-xs space-y-1' }, [
                h('p', '管理员: admin / 123456')
              ]),
              h('div', { class: 'flex gap-2 mt-3' }, [
                h('button', {
                  type: 'button',
                  onClick: () => {
                    loginForm.username = 'admin';
                    loginForm.password = '123456';
                  },
                  class: 'w-full px-3 py-2 bg-amber-100 text-amber-800 rounded-xl text-xs font-bold hover:bg-amber-200 transition-colors'
                }, '填充管理员')
              ])
            ]),
            h('p', { class: 'mb-2' }, '注册即表示您同意我们的服务条款和隐私政策'),
            h('p', [
              '遇到问题？',
              h('a', { 
                href: '#', 
                class: 'text-indigo-600 hover:text-indigo-500 font-medium ml-1',
                onClick: (e: Event) => e.preventDefault()
              }, '联系技术支持')
            ])
          ])
        ])
      ])
    ]);
  }
});