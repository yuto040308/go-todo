import api from '@/lib/api/client';

type HelloResponse = {
  message: string;
};

export const getHello = async (): Promise<HelloResponse> => {
  const response = await api.get<HelloResponse>('/hello');
  return response.data;
};
