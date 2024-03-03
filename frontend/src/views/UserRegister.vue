<script setup>
import baseUrl from '../../baseconfig';

const username = defineModel('username')
const email = defineModel('email')

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
    console.log(response)
}
</script>

<template>
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
                <input type="email" class="form-control" id="email" v-model="email" placeholder="Enter email" required>
            </div>
            <div class="btn-container">
                <button type="submit" class="btn btn-primary">Submit</button>
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