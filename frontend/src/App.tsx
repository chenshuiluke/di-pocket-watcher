import { useState } from "react";
import reactLogo from "./assets/react.svg";
import viteLogo from "/vite.svg";
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

function App() {
  const [opened, { toggle }] = useDisclosure();
  return (
    <MantineProvider>
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
            <Burger
              opened={opened}
              onClick={toggle}
              hiddenFrom="sm"
              size="sm"
            />
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
        <AppShell.Main>Main</AppShell.Main>
      </AppShell>
    </MantineProvider>
  );
}

export default App;
