import axios from 'axios'
import { Message } from 'element-ui'

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
      !ignoreMsg && !ignoreSuccessMsg && Message({type:'success',message:response.data.message || '操作成功'})
      return {
        ...response.data,
        success: true
      }
    case 300:
      Message({type:'error',message:response.data.message || '发生业务异常'})
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

export default function request () {

  function handleCatch(err) {
    console.log('请求抛出的异常信息', err)
    // message.error('请求异常')
    return { message: '请求异常' }
  }

  function handleResponse(response) {
    return response
  }

  function get(url,obj) {
    return axios.get(url, {
      ...obj
    })
    .then(handleResponse)
    .catch(handleCatch)
  }

  function post(url,obj) {
    return axios.post(url, {
      ...obj
    })
    .then(handleResponse)
    .catch(handleCatch)
  }

  return {
    get:get,
    post:post
  }

}

