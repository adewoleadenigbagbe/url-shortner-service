<script setup>
import HelloWorld from '@/components/HelloWorld.vue';
import { useRouter } from 'vue-router';
import baseUrl from '../../baseconfig';

const username = defineModel('username')
const email = defineModel('email')
const router = useRouter()

async function registerUser() {
    const requestOptions = {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            username: username.value,
            email: email.value
        })
    };

    const resource = "auth/register"
    let url = `${baseUrl}${resource}`
    const response = await fetch(url, requestOptions)
    if (response.ok && response.status === 200) {
        const data = await response.json()
        localStorage.setItem("apikey", data["apikey"])
        router.push({ name: 'login' })
    }

    console.log(response)
}

// function displayGrowlerError(){
//     console.log("hi")
// }
</script>

<template>
    <header>
        <div class="wrapper">
            <HelloWorld welcomeMsg="One more step" moreInfo="Experience full control over your short links" />
        </div>
    </header>
    <main>
        <div>
            <h1>Register Details</h1>
            <form @submit.prevent="registerUser">
                <div>
                    <label for="username" class="form-label mt-4">UserName</label>
                    <input type="text" class="form-control" id="username" v-model="username"
                        placeholder="Enter your username" required>
                </div>
                <div>
                    <label for="email" class="form-label mt-4">Email address</label>
                    <input type="email" class="form-control" id="email" v-model="email" placeholder="Enter email"
                        required>
                </div>
                <div class="btn-container">
                    <button type="submit" class="btn btn-primary">Register</button>
                </div>
            </form>
        </div>
    </main>
</template>

<style scoped>
header {
  line-height: 1.5;
  max-height: 100vh;
}

h1 {
    font-weight: 400;
}

.btn-container {
    margin-top: 20px;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
}
</style>