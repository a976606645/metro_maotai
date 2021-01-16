import request from '@/utils/request'

export function storeList(data) {
    return request({
        url: '/store/list',
        method: 'post',
        data
    })
}