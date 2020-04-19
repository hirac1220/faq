<template>
  <div class="container">
    <div>
      <h1 class="title">
        お家ですごそう
        Stay@Home
      </h1>
      <div>{{ text }}</div>
      <nuxt-link to="/faq">
        Click!
      </nuxt-link>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      text: ''
    }
  }
}
</script>

<style>
/* responsive */
* {
    box-sizing: border-box;
}
/* @media all and (max-width: 400px) { */
  .container {
    margin: 0 auto;
    min-height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
    text-align: center;
  }

  .title {
    font-family: 'Quicksand', 'Source Sans Pro', -apple-system, BlinkMacSystemFont,
      'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
    display: block;
    font-weight: 300;
    font-size: 45px;
    color: #35495e;
    letter-spacing: 1px;
  }

  .subtitle {
    font-weight: 300;
    font-size: 45px;
    color: #526488;
    word-spacing: 5px;
    padding-bottom: 15px;
  }

  .links {
    padding-top: 15px;
  }
/* } */
</style>
