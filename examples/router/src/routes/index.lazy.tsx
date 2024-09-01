import { createLazyFileRoute } from '@tanstack/react-router'

export const Route = createLazyFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <div className="px-10 py-4">
      <h1 className="text-xl leading-10 font-medium">Home</h1>
      <p className="py-2">
        This is a simple example of a fullstack app built with
        TanStack Router and TanStack Query in the frontend and Go templating
        with Vite-rendered assets in the backend.
      </p>
      <ul className="list-disc list-outside">
        <li>Use the navigation links above to explore the app.</li>
        <li>The Vite icon is rendered from the application's public folder.</li>
        <li>The Home link renders the home page (surprise!).</li>
        <li>
          The About page is a different page registered on TanStack Router.
          You can click on it client-side, but can also do a
          browser refresh to render it from the server.
        </li>
        <li>
          The Dad Joke calls an endpoint on the server to fetch a random dad joke.
          Refresh the page to see a new joke. It illustrates how frontend and
          backend can work together in a single codebase.
        </li>
      </ul>
    </div>
  )
}
