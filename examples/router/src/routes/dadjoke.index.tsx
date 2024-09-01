import { useSuspenseQuery } from '@tanstack/react-query';
import { createFileRoute, Link } from '@tanstack/react-router';
import { zodSearchValidator } from '@tanstack/router-zod-adapter';
import { z } from 'zod';
import { dadJokesQueryOptions } from '../api/jokes';

const dadJokeSearchSchema = z.object({
  t: z.number().int().default(Date.now()),
})

export const Route = createFileRoute('/dadjoke/')({
  component: DadJokeComponent,
  validateSearch: zodSearchValidator(dadJokeSearchSchema),
})

function DadJokeComponent() {
  const { t } = Route.useSearch();
  const { data: dadjoke } = useSuspenseQuery(dadJokesQueryOptions.random(t));

  return (
    <div className="px-10 py-4">
      <h1 className="text-xl leading-10 font-medium">Dad Joke</h1>
      <div>{dadjoke.joke}</div>
      <Link
        to="/dadjoke"
        search={{
          t: Date.now(),
        }}
        replace
        className="mt-2 relative inline-flex items-center justify-center whitespace-nowrap rounded-md border px-3 py-2 text-center text-sm font-medium shadow-sm transition-all duration-100 ease-in-out disabled:pointer-events-none disabled:shadow-none outline outline-offset-2 outline-0 focus-visible:outline-2 outline-blue-500 dark:outline-blue-500 border-gray-300 dark:border-gray-800 text-gray-900 dark:text-gray-50 bg-white dark:bg-gray-950 hover:bg-gray-50 dark:hover:bg-gray-900/60 disabled:text-gray-400 disabled:dark:text-gray-600">
          Refresh
      </Link>
    </div>
  )
}
