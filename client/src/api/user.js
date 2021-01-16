import request from '@/utils/request'

export function userList(data) {
  return request({
    url: '/user/list',
    method: 'post',
    data
  })
}

export function sendSms(data) {
  return request({
    url: '/user/sendSms',
    method: 'post',
    data
  })
}

export function userLogin(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

export function userAdd(data) {
  return request({
    url: '/user/add',
    method: 'post',
    data
  })
}

export function setStore(data) {
  return request({
    url: '/user/setStore',
    method: 'post',
    data
  })
}