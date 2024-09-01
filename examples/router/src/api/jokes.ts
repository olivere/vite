import { queryOptions } from '@tanstack/react-query';
import axios from 'redaxios';

export type DadJoke = {
  id: string;
  joke: string;
  status: number;
}

export const fetchRandomDadJoke = async (delay?: number): Promise<DadJoke> => {
  if (delay) {
    await new Promise((resolve) => setTimeout(resolve, delay));
  }
  const response = await axios.get<DadJoke>('http://localhost:8080/api/dadjoke', {
    headers: {
      Accept: 'application/json',
    },
  });

  return response.data;
};

export const dadJokesQueryOptions = {
  random: (random: number, delay?: number) => queryOptions({
    queryKey: ['dadjokes', random],
    queryFn: () => fetchRandomDadJoke(delay),
  }),
};
