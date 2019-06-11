import axios from 'axios'
import { Message } from 'element-ui'
import router from '@/router'

//请求队列，防止重复发送
let xhrQueue = []
const delXhrQueue = (url) => {
  const idx = xhrQueue.findIndex((i) => {
    return url.search(i) > -1
  })
  xhrQueue.splice(idx, 1)
}

axios.interceptors.request.use((req) => {
  if (!req.ignoreQueen) {
    if (xhrQueue.includes(req.url)) {
      return null
    }
    xhrQueue.push(req.url)
  }

  let token = localStorage.getItem('login_token');
  if ( token!=null && token!='' ) {
    req.headers.Authorization = `Bearer ${token}`
  }

  return req
})

axios.interceptors.response.use((response) => {
  const {
    // ContentType,
    url,
    ignoreMsg = false, // true 拒绝公共显示信息
    ignoreSuccessMsg = false, // true 拒绝成功提示信息
    allResponseData
  } = response.config
  // 去除拦截
  delXhrQueue(url)
  if (allResponseData) {
    return {
      ...response
    }
  }
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
      !ignoreMsg && Message({type:'error',message:'未定义的异常'})
      break
  }
  return response.data
}, (error) => {
  return Promise.reject(error)
})