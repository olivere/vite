import { createLazyFileRoute, Outlet } from '@tanstack/react-router';

export const Route = createLazyFileRoute('/dadjoke')({
  component: Outlet,
})
