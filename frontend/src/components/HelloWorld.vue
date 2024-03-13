<script setup>
defineProps({
  welcomeMsg: {
    type: String,
    required: true
  },
  moreInfo: {
    type: String,
    required: true
  },
})
</script>

<template>
  <img alt="Vue logo" class="logo" src="@/assets/logo.svg" width="125" height="125" />
  <div class="greetings">
    <h1 class="green">{{ welcomeMsg }}</h1>
    <h3>{{ moreInfo }}</h3>
  </div>
</template>

<style scoped>
h1 {
  font-weight: 500;
  font-size: 2.4rem;
  position: relative;
  top: -10px;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}


.greetings h1,
.greetings h3 {
  text-align: center;
}

@media (min-width: 1024px) {
  .greetings h1{
    text-align: left;
  }
}
</style>
