import axios from 'axios'
import { Message } from 'element-ui'
import router from '@/router'


axios.interceptors.request.use((req) => {
  let token = localStorage.getItem('login_token');
  if ( token!=null && token!='' ) {
    req.headers.Authorization = `Bearer ${token}`
  }

  return req
}, (err) => {
  return Promise.reject(err)
})

axios.interceptors.response.use((response) => {
  switch (response.data.code) {
    case 0:
    case 200:
      if (response.data.message){
        Message({type:'success',message:response.data.message})
      }
      return {
        ...response.data,
        success: true
      }
    case 300:
      Message({type:'error',message:response.data.message || '发生业务异常'})
      break
    case 400:
    case 401:
      Message({type:'error',message:response.data.message || '登录状态无效,请重新登录'})
      router.push("/login")
      break
    case 403:
      Message({type:'error',message:response.data.message || '无权访问'})
      router.push("/403")
      break
    case 500:
      Message({type:'error',message:response.data.message || '服务器端异常'})
      break
    default:
      break
  }
  return response
}, (error) => {
  return Promise.reject(error)
})