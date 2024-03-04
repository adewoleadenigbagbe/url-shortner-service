<script setup>
import baseUrl from '../../baseconfig';

const email = defineModel('email')

async function loginUser() {
    const requestOptions = {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            email: email.value
        })
    };

    const resource = "auth/sign-in"
    let url = `${baseUrl}${resource}`
    const response = await fetch(url, requestOptions)
    console.log(response)
}
</script>

<template>
    <div>
        <h1>Log In</h1>
        <form @submit.prevent="loginUser">
            <div>
                <label for="email" class="form-label mt-4">Email address</label>
                <input type="email" class="form-control" id="email" v-model="email" placeholder="Enter email" required>
            </div>
            <div class="btn-container">
                <button type="submit" class="btn btn-primary">LogIn</button>
            </div>
        </form>
    </div>
</template>

<style scoped>
h1 {
    font-weight: 400;
}

.btn-container {
    margin-top: 20px;
}
</style>