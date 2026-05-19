'use client';

import { useQuery } from '@tanstack/react-query';
import { getHello } from '@/api/hello';

export default function Home() {
  const { data, isLoading } = useQuery({
    queryKey: ['hello'],
    queryFn: getHello,
  });

  if (isLoading || !data) return <p>Loading...</p>;

  return <p>{data.message}</p>;
}
