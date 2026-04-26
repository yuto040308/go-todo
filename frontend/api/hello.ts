import api from "@/lib/api"

type HelloResponse = {
    message: string
}

export const getHello = async (): Promise<HelloResponse> => {
    const response = await api.get('/hello')
    return response.data
}