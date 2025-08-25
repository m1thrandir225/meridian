<script setup lang="ts">
import ChannelsSidebar from '@/components/ChannelsSidebar.vue'
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Plus, MessageSquare, Users, Zap, Menu, ArrowRight } from 'lucide-vue-next'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useChannelStore } from '@/stores/channel'
import { onMounted } from 'vue'
import Logo from '@/components/LogoWord.vue'

const router = useRouter()
const channelStore = useChannelStore()
const isChannelsSidebarOpen = ref(true)
const showCreateDialog = ref(false)

const toggleChannelsSidebar = () => {
  isChannelsSidebarOpen.value = !isChannelsSidebarOpen.value
}

onMounted(async () => {
  if (!channelStore.channels.length) {
    await channelStore.fetchChannels()
  }
})

const handleCreateChannel = () => {
  showCreateDialog.value = true
}

const features = [
  {
    icon: MessageSquare,
    title: 'Real-time Messaging',
    description:
      'Send and receive messages instantly with real-time updates across all your devices.',
  },
  {
    icon: Users,
    title: 'Team Collaboration',
    description:
      'Create channels for different topics and invite team members to collaborate effectively.',
  },
  {
    icon: Zap,
    title: 'Bot Integrations',
    description:
      'Enhance your workflow with custom bots and integrations for automation and productivity.',
  },
]
</script>

<template>
  <div class="flex h-screen">
    <SidebarProvider>
      <Transition
        name="slide-left"
        enter-active-class="transition-all duration-300 ease-out"
        leave-active-class="transition-all duration-300 ease-in"
        enter-from-class="-translate-x-full opacity-0"
        enter-to-class="translate-x-0 opacity-100"
        leave-from-class="translate-x-0 opacity-100"
        leave-to-class="-translate-x-full opacity-0"
      >
        <ChannelsSidebar
          v-if="isChannelsSidebarOpen"
          v-model:show-create-dialog="showCreateDialog"
        />
      </Transition>
      <SidebarInset>
        <!-- Welcome Content Area -->
        <div class="flex flex-col h-full bg-background">
          <!-- Header -->
          <header
            class="flex items-center justify-between px-6 py-4 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
          >
            <div class="flex items-center gap-3">
              <Button
                variant="default"
                size="icon"
                @click="toggleChannelsSidebar"
                :class="isChannelsSidebarOpen ? '' : 'bg-accent'"
              >
                <Menu class="h-4 w-4" />
              </Button>
              <h1 class="text-xl font-semibold">Welcome to Meridian</h1>
            </div>
          </header>

          <!-- Main Content -->
          <div class="flex-1 overflow-y-auto p-6">
            <div class="max-w-4xl mx-auto">
              <!-- Welcome Section -->
              <div class="text-center mb-12">
                <div class="mb-6">
                  <Logo :height="80" class="mx-auto mb-4" />
                </div>
                <h2 class="text-3xl font-bold mb-4">Welcome to Meridian</h2>
                <p class="text-lg text-muted-foreground max-w-2xl mx-auto mb-8">
                  Your modern messaging platform for seamless team communication and collaboration.
                  Get started by creating your first channel or exploring existing ones.
                </p>

                <!-- Quick Actions -->
                <div class="flex items-center justify-center gap-4 mb-12">
                  <Button variant="outline" size="lg" @click="handleCreateChannel" class="gap-2">
                    <Plus class="h-4 w-4" />
                    Create Channel
                  </Button>
                  <Button
                    variant="outline"
                    size="lg"
                    @click="router.push('/bot-management')"
                    class="gap-2"
                  >
                    <Zap class="h-4 w-4" />
                    Manage Bots
                  </Button>
                </div>
              </div>

              <!-- Features Grid -->
              <div class="grid md:grid-cols-3 gap-6 mb-12">
                <Card v-for="feature in features" :key="feature.title" class="text-center">
                  <CardHeader>
                    <div
                      class="mx-auto mb-4 w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center"
                    >
                      <component :is="feature.icon" class="h-6 w-6 text-primary" />
                    </div>
                    <CardTitle class="text-lg">{{ feature.title }}</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <CardDescription>{{ feature.description }}</CardDescription>
                  </CardContent>
                </Card>
              </div>

              <!-- Channel Status -->
              <Card v-if="channelStore.channels.length > 0">
                <CardHeader>
                  <CardTitle>Your Channels</CardTitle>
                  <CardDescription>
                    You have {{ channelStore.channels.length }} channel{{
                      channelStore.channels.length !== 1 ? 's' : ''
                    }}
                    available.
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div class="grid gap-2">
                    <div
                      v-for="channel in channelStore.channels.slice(0, 3)"
                      :key="channel.id"
                      class="flex items-center justify-between p-3 rounded-lg border hover:bg-accent/50 transition-colors cursor-pointer"
                      @click="router.push(`/channel/${channel.id}`)"
                    >
                      <div class="flex items-center gap-3">
                        <div
                          class="w-8 h-8 bg-primary/10 rounded-lg flex items-center justify-center"
                        >
                          <MessageSquare class="h-4 w-4 text-primary" />
                        </div>
                        <div>
                          <p class="font-medium">{{ channel.name }}</p>
                          <p class="text-sm text-muted-foreground">{{ channel.topic }}</p>
                        </div>
                      </div>
                      <Button as-child variant="outline" size="icon">
                        <RouterLink :to="`/channel/${channel.id}`">
                          <ArrowRight class="h-4 w-4" />
                        </RouterLink>
                      </Button>
                    </div>
                    <div v-if="channelStore.channels.length > 3" class="text-center pt-2">
                      <Button variant="outline" @click="router.push('/')">
                        View All Channels
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <!-- Empty State -->
              <Card v-else>
                <CardHeader>
                  <CardTitle>No Channels Yet</CardTitle>
                  <CardDescription>
                    Create your first channel to start collaborating with your team.
                  </CardDescription>
                </CardHeader>
                <CardContent class="text-center">
                  <div class="mb-4">
                    <MessageSquare class="h-12 w-12 text-muted-foreground mx-auto" />
                  </div>
                  <Button @click="handleCreateChannel" class="gap-2">
                    <Plus class="h-4 w-4" />
                    Create Your First Channel
                  </Button>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  </div>
</template>
