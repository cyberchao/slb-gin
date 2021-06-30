import service from '@/utils/request'

export const cross_parse = (data) => {
    return service({
        url: "/api/cross_parse",
        method: 'post',
        data:data,
    })
}