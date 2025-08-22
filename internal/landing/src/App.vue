<script setup lang="ts">
// Component imports
import FeatureCard from './components/FeatureCard.vue'
import TechCard from './components/TechCard.vue'
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- Navigation -->
    <nav class="fixed top-0 w-full bg-white/80 backdrop-blur-md border-b border-slate-200 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <h1
                class="text-2xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent"
              >
                Meridian
              </h1>
            </div>
          </div>
          <div class="hidden md:block">
            <div class="ml-10 flex items-baseline space-x-4">
              <a
                href="#features"
                class="text-slate-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              >
                Features
              </a>
              <a
                href="#about"
                class="text-slate-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              >
                About
              </a>
              <a
                href="#tech"
                class="text-slate-600 hover:text-blue-600 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              >
                Technology
              </a>
              <a
                href="http://app.localhost"
                class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-full text-sm font-medium transition-colors"
              >
                Launch App
              </a>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="pt-64 pb-36 px-4 sm:px-6 lg:px-8">
      <div class="max-w-7xl mx-auto">
        <div class="text-center">
          <h1 class="text-4xl sm:text-5xl lg:text-6xl font-bold text-slate-900 mb-6">
            Modern
            <span
              class="bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent"
            >
              Communication
            </span>
            Platform
          </h1>
          <p class="text-xl text-slate-600 max-w-3xl mx-auto mb-8 leading-relaxed">
            Connect, collaborate, and communicate seamlessly with our powerful messaging platform.
            Built for teams that demand reliability, security, and performance.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <a
              href="http://app.localhost"
              class="bg-blue-600 hover:bg-blue-700 text-white px-8 py-3 rounded-full text-lg font-semibold transition-colors shadow-lg hover:shadow-xl"
            >
              Launch App
            </a>
            <a
              href="#features"
              class="border-2 border-indigo-600 text-indigo-600 hover:bg-indigo-600 hover:text-white px-8 py-3 rounded-full text-lg font-semibold transition-colors"
            >
              Learn More
            </a>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="features" class="py-16 px-4 sm:px-6 lg:px-8 bg-white">
      <div class="max-w-7xl mx-auto">
        <div class="text-center mb-16">
          <h2 class="text-3xl sm:text-4xl font-bold text-slate-900 mb-4">Why Choose Meridian?</h2>
          <p class="text-xl text-slate-600 max-w-2xl mx-auto">
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
          <FeatureCard
            icon=""
            title="Modern Stack"
            description="Built with Go, Vue.js, and modern web technologies."
          />
        </div>
      </div>
    </section>

    <!-- Demo Section -->
    <section class="py-16 px-4 sm:px-6 lg:px-8 bg-gradient-to-br from-slate-50 to-blue-50">
      <div class="max-w-7xl mx-auto">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div>
            <h2 class="text-3xl sm:text-4xl font-bold text-slate-900 mb-6">
              Experience the Future of Communication
            </h2>
            <p class="text-xl text-slate-600 mb-8 leading-relaxed">
              Our platform combines the best of modern messaging with enterprise-grade security and
              performance. See how teams can collaborate more effectively with real-time
              communication.
            </p>
            <div class="space-y-4">
              <div class="flex items-center">
                <div
                  class="w-6 h-6 bg-green-500 rounded-full flex items-center justify-center mr-3"
                >
                  <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path
                      fill-rule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <span class="text-slate-700">Real-time message synchronization</span>
              </div>
              <div class="flex items-center">
                <div
                  class="w-6 h-6 bg-green-500 rounded-full flex items-center justify-center mr-3"
                >
                  <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path
                      fill-rule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <span class="text-slate-700">Secure authentication and authorization</span>
              </div>
              <div class="flex items-center">
                <div
                  class="w-6 h-6 bg-green-500 rounded-full flex items-center justify-center mr-3"
                >
                  <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                    <path
                      fill-rule="evenodd"
                      d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <span class="text-slate-700">Integration with external services</span>
              </div>
            </div>
          </div>
          <div class="relative">
            <div
              class="bg-white rounded-2xl shadow-2xl p-6 transform rotate-1 hover:rotate-0 transition-transform duration-300"
            >
              <div class="flex items-center mb-4">
                <div class="w-3 h-3 bg-red-500 rounded-full mr-2"></div>
                <div class="w-3 h-3 bg-yellow-500 rounded-full mr-2"></div>
                <div class="w-3 h-3 bg-green-500 rounded-full mr-4"></div>
                <div class="text-sm text-slate-500">Meridian Chat</div>
              </div>
              <div class="space-y-3">
                <div class="flex items-start space-x-3">
                  <div
                    class="w-8 h-8 bg-blue-500 rounded-full flex items-center justify-center text-white text-sm font-semibold"
                  >
                    A
                  </div>
                  <div class="flex-1">
                    <div class="text-sm font-semibold text-slate-900">Alice</div>
                    <div class="bg-blue-100 rounded-lg p-3 text-slate-600 text-sm">
                      Welcome to Meridian! 🚀
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
                    <div class="text-sm font-semibold text-slate-900">Bob</div>
                    <div class="bg-green-100 rounded-lg text-slate-600 p-3 text-sm">
                      This platform looks amazing! The real-time features are incredible.
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
                    <div class="text-sm font-semibold text-slate-900">Charlie</div>
                    <div class="bg-purple-100 rounded-lg text-slate-600 p-3 text-sm">
                      Agreed! The integration capabilities are game-changing.
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Technology Stack Section -->
    <section id="tech" class="py-16 px-4 sm:px-6 lg:px-8 bg-white">
      <div class="max-w-7xl mx-auto">
        <div class="text-center mb-16">
          <h2 class="text-3xl sm:text-4xl font-bold text-slate-900 mb-4">
            Built with Modern Technologies
          </h2>
          <p class="text-xl text-slate-600 max-w-2xl mx-auto">
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
    <section
      id="about"
      class="py-16 px-4 sm:px-6 lg:px-8 bg-gradient-to-br from-slate-50 to-blue-50"
    >
      <div class="max-w-7xl mx-auto">
        <div class="text-center mb-16">
          <h2 class="text-3xl sm:text-4xl font-bold text-slate-900 mb-4">About Meridian</h2>
          <p class="text-xl text-slate-600 max-w-3xl mx-auto leading-relaxed">
            Meridian is a modern communication platform built with microservices architecture. It
            provides real-time messaging, secure authentication, and seamless integration
            capabilities for teams and organizations that demand reliability and performance.
          </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="text-center">
            <div
              class="w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <svg
                class="w-8 h-8 text-blue-600"
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
            <h3 class="text-xl font-semibold text-slate-900 mb-2">High Performance</h3>
            <p class="text-slate-600">Built with Go for exceptional speed and efficiency.</p>
          </div>
          <div class="text-center">
            <div
              class="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <svg
                class="w-8 h-8 text-green-600"
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
            <h3 class="text-xl font-semibold text-slate-900 mb-2">Secure by Design</h3>
            <p class="text-slate-600">Enterprise-grade security with end-to-end encryption.</p>
          </div>
          <div class="text-center">
            <div
              class="w-16 h-16 bg-purple-100 rounded-full flex items-center justify-center mx-auto mb-4"
            >
              <svg
                class="w-8 h-8 text-purple-600"
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
            <h3 class="text-xl font-semibold text-slate-900 mb-2">Developer Friendly</h3>
            <p class="text-slate-600">Comprehensive API and integration capabilities.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-slate-900 text-white py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-7xl mx-auto">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div>
            <h3
              class="text-2xl font-bold bg-gradient-to-r from-blue-400 to-indigo-400 bg-clip-text text-transparent mb-4"
            >
              Meridian
            </h3>
            <p class="text-slate-400">
              Modern communication platform for teams that demand reliability and performance.
            </p>
          </div>
          <div>
            <h4 class="text-lg font-semibold mb-4">Product</h4>
            <ul class="space-y-2 text-slate-400">
              <li><a href="#features" class="hover:text-white transition-colors">Features</a></li>
              <li><a href="#tech" class="hover:text-white transition-colors">Technology</a></li>
              <li>
                <a href="http://app.localhost" class="hover:text-white transition-colors"
                  >Launch App</a
                >
              </li>
            </ul>
          </div>
          <div>
            <h4 class="text-lg font-semibold mb-4">Connect</h4>
            <ul class="space-y-2 text-slate-400">
              <li><a href="#" class="hover:text-white transition-colors">Documentation</a></li>
              <li><a href="#" class="hover:text-white transition-colors">Support</a></li>
              <li><a href="#" class="hover:text-white transition-colors">GitHub</a></li>
            </ul>
          </div>
        </div>
        <div class="border-t border-slate-800 mt-8 pt-8 text-center text-slate-400">
          <p>&copy; 2025 Meridian. Built by Sebastijan Zindl.</p>
        </div>
      </div>
    </footer>
  </div>
</template>
