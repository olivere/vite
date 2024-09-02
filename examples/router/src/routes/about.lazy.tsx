import { createLazyFileRoute } from '@tanstack/react-router';

export const Route = createLazyFileRoute('/about')({
  component: About,
})

function About() {
  return (
    <div className="px-10 py-4">
      <h1 className="text-xl leading-10 font-medium">About</h1>
      <p>Just another page rendered from the client-side.</p>
    </div>
  );
}
