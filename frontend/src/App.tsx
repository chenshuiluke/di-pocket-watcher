import "./App.css";
import "@mantine/core/styles.css";

import {
  AppShell,
  Burger,
  Group,
  MantineProvider,
  Skeleton,
  Title,
} from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import { Login } from "./components/Login";
import { QueryProvider } from "./lib/providers/query-provider";
import { useCurrentUser } from "./features/auth/api/auth";

function AppContent() {
  const [opened, { toggle }] = useDisclosure();
  const { data: user, isLoading } = useCurrentUser();

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{
        width: 300,
        breakpoint: "sm",
        collapsed: { mobile: !opened },
      }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md">
          <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
          <Title order={1}>Di Pocket Watcher</Title>
        </Group>
      </AppShell.Header>
      <AppShell.Navbar p="md">
        Navbar
        {Array(15)
          .fill(0)
          .map((_, index) => (
            <Skeleton key={index} h={28} mt="sm" animate={false} />
          ))}
      </AppShell.Navbar>
      <AppShell.Main>
        {!user ? <Login /> : <div>Welcome, {user.email}!</div>}
      </AppShell.Main>
    </AppShell>
  );
}

function App() {
  return (
    <QueryProvider>
      <MantineProvider>
        <AppContent />
      </MantineProvider>
    </QueryProvider>
  );
}

export default App;
