<script setup lang="ts">
// Component imports
import FeatureCard from './components/FeatureCard.vue'
import TechCard from './components/TechCard.vue'
import LogoWord from './components/LogoWord.vue'
import AppearanceMenu from './components/AppearanceMenu.vue'
import { useAppearanceStore } from '@/stores/appearance'
import { onMounted } from 'vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { useHead } from '@unhead/vue'

const appearanceStore = useAppearanceStore()

onMounted(() => {
  appearanceStore.initializeAppearance()
})

useHead({
  title: 'Meridian',
  meta: [
    {
      name: 'description',
      content:
        'Meridian is a modern communication platform built with microservices architecture. It provides real-time messaging, secure authentication, customizable themes, and seamless integration capabilities for teams and organizations that demand reliability, performance, and a personalized experience.',
    },
  ],
  link: [
    {
      rel: 'icon',
      type: 'image/png',
      href: '/favicon.png',
    },
  ],
})
</script>

<template>
  <div class="min-h-screen bg-background">
    <!-- Navigation -->
    <nav class="fixed top-0 w-full bg-background/80 backdrop-blur-md border-b border-border z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <LogoWord :height="32" />
            </div>
          </div>
          <div class="hidden md:flex items-center space-x-4">
            <a
              href="#features"
              class="text-muted-foreground hover:text-foreground px-3 py-2 rounded-md text-sm font-medium transition-colors"
            >
              Features
            </a>
            <a
              href="#about"
              class="text-muted-foreground hover:text-foreground px-3 py-2 rounded-md text-sm font-medium transition-colors"
            >
              About
            </a>
            <a
              href="#tech"
              class="text-muted-foreground hover:text-foreground px-3 py-2 rounded-md text-sm font-medium transition-colors"
            >
              Technology
            </a>
            <AppearanceMenu />
            <Button as-child>
              <a href="http://app.localhost">Launch App</a>
            </Button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="pt-64 pb-36 px-4 sm:px-6 lg:px-8">
      <div class="max-w-7xl mx-auto">
        <div class="text-center">
          <div class="mb-8">
            <LogoWord :height="120" class="mx-auto" />
          </div>
          <h1 class="text-4xl sm:text-5xl lg:text-6xl font-bold text-foreground mb-6">
            Your Modern
            <span class="text-primary">Messaging Platform</span>
          </h1>
          <p class="text-xl text-muted-foreground max-w-3xl mx-auto mb-8 leading-relaxed">
            Meridian is a powerful, real-time messaging platform built for modern teams. Experience
            seamless communication with enterprise-grade security, customizable themes, and powerful
            bot integrations. Connect, collaborate, and communicate like never before.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <Button size="lg" as-child>
              <a href="http://app.localhost">Launch App</a>
            </Button>
            <Button variant="outline" size="lg" as-child>
              <a href="#features">Learn More</a>
            </Button>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="py-16 px-4 sm:px-6 lg:px-8 bg-card">
      <div class="max-w-7xl mx-auto">
        <div class="text-center mb-16">
          <h2 class="text-3xl sm:text-4xl font-bold text-foreground mb-4">Why Choose Meridian?</h2>
          <p class="text-xl text-muted-foreground max-w-2xl mx-auto">
            Built with modern technologies and designed for the future of team communication.
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          <FeatureCard
            icon="💬"
            title="Real-time Messaging"
            description="Instant messaging with real-time updates, typing indicators, and message reactions."
          />
          <FeatureCard
            icon="🎨"
            title="Customizable Themes"
            description="Personalize your experience with multiple themes and accent colors to match your style."
          />
          <FeatureCard
            icon="🔒"
            title="Secure & Private"
            description="End-to-end encryption, secure authentication, and privacy-first design."
          />
          <FeatureCard
            icon="🤖"
            title="Bot Integration"
            description="Integrate with your favorite tools and services through our powerful API."
          />
          <FeatureCard
            icon="🌐"
            title="Cross-platform"
            description="Access from anywhere with our responsive web app that works on all devices."
          />
          <FeatureCard
            icon="⚡"
            title="High Performance"
            description="Built with microservices architecture for scalability and reliability."
          />
          <FeatureCard
            icon="🔄"
            title="Event-driven"
            description="Real-time updates powered by Apache Kafka for seamless communication."
          />
          <FeatureCard
            icon="🗄️"
            title="Persistent Storage"
            description="Reliable PostgreSQL database with Redis caching for optimal performance."
          />
        </div>
      </div>
    </section>

    <!-- Demo Section -->
    <section class="py-16 px-4 sm:px-6 lg:px-8 bg-background">
      <div class="max-w-7xl mx-auto">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <h2 class="text-3xl sm:text-4xl font-bold text-foreground mb-6">
              Experience the Future of Communication
            </h2>
            <p class="text-xl text-muted-foreground mb-8 leading-relaxed">
              Our platform combines the best of modern messaging with enterprise-grade security and
              performance. See how teams can collaborate more effectively with real-time
              communication, customizable themes, and powerful integrations.
            </p>
            <div class="space-y-4">
              <div class="flex items-center">
                <div class="w-6 h-6 bg-primary rounded-full flex items-center justify-center mr-3">
                  <svg
                    class="w-4 h-4 text-primary-foreground"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <span class="text-foreground">Real-time message synchronization</span>
              </div>
              <div class="flex items-center">
                <div class="w-6 h-6 bg-primary rounded-full flex items-center justify-center mr-3">
                  <svg
                    class="w-4 h-4 text-primary-foreground"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <span class="text-foreground">Customizable themes and appearance</span>
              </div>
              <div class="flex items-center">
                <div class="w-6 h-6 bg-primary rounded-full flex items-center justify-center mr-3">
                  <svg
                    class="w-4 h-4 text-primary-foreground"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <span class="text-foreground">Secure authentication and authorization</span>
              </div>
              <div class="flex items-center">
                <div class="w-6 h-6 bg-primary rounded-full flex items-center justify-center mr-3">
                  <svg
                    class="w-4 h-4 text-primary-foreground"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <span class="text-foreground">Integration with external services</span>
              </div>
            </div>
          </div>
          <div class="relative">
            <Card class="transform rotate-1 hover:rotate-0 transition-transform duration-300">
              <CardHeader class="pb-3">
                <div class="flex items-center justify-between">
                  <div class="flex items-center space-x-2">
                    <div class="w-3 h-3 bg-red-500 rounded-full"></div>
                    <div class="w-3 h-3 bg-yellow-500 rounded-full"></div>
                    <div class="w-3 h-3 bg-green-500 rounded-full"></div>
                  </div>
                  <span class="text-sm text-muted-foreground">Meridian Chat</span>
                </div>
              </CardHeader>
              <CardContent class="space-y-3">
                <div class="flex items-start space-x-3">
                  <div
                    class="w-8 h-8 bg-primary rounded-full flex items-center justify-center text-primary-foreground text-sm font-semibold"
                  >
                    A
                  </div>
                  <div class="flex-1">
                    <div class="text-sm font-semibold text-foreground">Alice</div>
                    <div class="bg-muted rounded-lg p-3 text-foreground text-sm">
                      Welcome to Meridian! 🚀 The new theme looks amazing!
                    </div>
                  </div>
                </div>
                <div class="flex items-start space-x-3">
                  <div
                    class="w-8 h-8 bg-green-500 rounded-full flex items-center justify-center text-white text-sm font-semibold"
                  >
                    B
                  </div>
                  <div class="flex-1">
                    <div class="text-sm font-semibold text-foreground">Bob</div>
                    <div class="bg-muted rounded-lg text-foreground p-3 text-sm">
                      This platform is incredible! The real-time features and customization options
                      are perfect.
                    </div>
                  </div>
                </div>
                <div class="flex items-start space-x-3">
                  <div
                    class="w-8 h-8 bg-purple-500 rounded-full flex items-center justify-center text-white text-sm font-semibold"
                  >
                    C
                  </div>
                  <div class="flex-1">
                    <div class="text-sm font-semibold text-foreground">Charlie</div>
                    <div class="bg-muted rounded-lg text-foreground p-3 text-sm">
                      Agreed! The bot integrations and security features are game-changing.
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </section>

    <!-- Technology Stack Section -->
    <section id="tech" class="py-16 px-4 sm:px-6 lg:px-8 bg-card">
      <div class="max-w-7xl mx-auto">
        <div class="text-center mb-16">
          <h2 class="text-3xl sm:text-4xl font-bold text-foreground mb-4">
            Built with Modern Technologies
          </h2>
          <p class="text-xl text-muted-foreground max-w-2xl mx-auto">
            Our platform leverages cutting-edge technologies for optimal performance and
            scalability.
          </p>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-8">
          <TechCard name="Go" color="bg-blue-200" />
          <TechCard name="VueJS" color="bg-green-200" />
          <TechCard name="PostgreSQL" color="bg-indigo-200" />
          <TechCard name="Redis" color="bg-red-200" />
          <TechCard name="Kafka" color="bg-orange-200" />
          <TechCard name="Docker" color="bg-indigo-200" />
        </div>
      </div>
    </section>

    <!-- About Section -->
    <section id="about" class="py-16 px-4 sm:px-6 lg:px-8 bg-background">
      <div class="max-w-7xl mx-auto">
        <div class="text-center mb-16">
          <h2 class="text-3xl sm:text-4xl font-bold text-foreground mb-4">About Meridian</h2>
          <p class="text-xl text-muted-foreground max-w-3xl mx-auto leading-relaxed">
            Meridian is a modern communication platform built with microservices architecture. It
            provides real-time messaging, secure authentication, customizable themes, and seamless
            integration capabilities for teams and organizations that demand reliability,
            performance, and a personalized experience.
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <Card class="text-center">
            <CardHeader>
              <div
                class="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4"
              >
                <svg
                  class="w-8 h-8 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M13 10V3L4 14h7v7l9-11h-7z"
                  />
                </svg>
              </div>
              <CardTitle>High Performance</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>Built with Go for exceptional speed and efficiency.</CardDescription>
            </CardContent>
          </Card>
          <Card class="text-center">
            <CardHeader>
              <div
                class="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4"
              >
                <svg
                  class="w-8 h-8 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
                  />
                </svg>
              </div>
              <CardTitle>Secure by Design</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription
                >Enterprise-grade security with end-to-end encryption.</CardDescription
              >
            </CardContent>
          </Card>
          <Card class="text-center">
            <CardHeader>
              <div
                class="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-4"
              >
                <svg
                  class="w-8 h-8 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                  />
                </svg>
              </div>
              <CardTitle>Developer Friendly</CardTitle>
            </CardHeader>
            <CardContent>
              <CardDescription>Comprehensive API and integration capabilities.</CardDescription>
            </CardContent>
          </Card>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-card border-t border-border py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-7xl mx-auto">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div>
            <LogoWord :height="40" class="mb-4" />
            <p class="text-muted-foreground">
              Modern communication platform for teams that demand reliability, performance, and
              customization.
            </p>
          </div>
          <div>
            <h4 class="text-lg font-semibold text-foreground mb-4">Product</h4>
            <ul class="space-y-2 text-muted-foreground">
              <li>
                <a href="#features" class="hover:text-foreground transition-colors">Features</a>
              </li>
              <li>
                <a href="#tech" class="hover:text-foreground transition-colors">Technology</a>
              </li>
              <li>
                <a href="http://app.localhost" class="hover:text-foreground transition-colors"
                  >Launch App</a
                >
              </li>
            </ul>
          </div>
          <div>
            <h4 class="text-lg font-semibold text-foreground mb-4">Connect</h4>
            <ul class="space-y-2 text-muted-foreground">
              <li><a href="#" class="hover:text-foreground transition-colors">Documentation</a></li>
              <li><a href="#" class="hover:text-foreground transition-colors">Support</a></li>
              <li><a href="#" class="hover:text-foreground transition-colors">GitHub</a></li>
            </ul>
          </div>
        </div>
        <Separator class="my-8" />
        <div class="text-center text-muted-foreground">
          <p>&copy; 2025 Meridian. Built by Sebastijan Zindl.</p>
        </div>
      </div>
    </footer>
  </div>
</template>
