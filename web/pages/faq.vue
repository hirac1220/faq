<template>
  <div class="container">
    <div>
      <div>
        <h2>どんな気分ですか？</h2>

        <div style="margin-bottom: 16px;">
          <input v-model="newFaq" type="text">
          <button v-on:click="add">
            Check!
          </button>
        </div>

        <ul>
          <li
            v-for="faq in faqs"
            v-bind:key="faq.ID"
            style="margin-bottom: 8px; text-align:left;"
          >
            <span
              key="default"
            >
              {{ faq.Text }}
            </span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data () {
    return {
      newFaq: '',
      faqs: [
        // APIから受け取るデータは、頭文字が大文字
        // { ID: 1, Text: 'Shopping', editable: false, newText: '' },
      ]
    }
  },
  // asyncDataは、ページ遷移前にAPIからデータを取得し、dataに反映します。
  // 反映したいdataのオブジェクトを返り値にします。
  async asyncData ({ $axios, error }) {
    let faqs = []

    // APIにGETメソッドでリクエストを送り、faqsを取得します。
    try {
      const { data } = await $axios.get('/faq')
      faqs = data
    } catch (err) {
      error({
        statusCode: err.response.status,
        message: err.response.statusText
      })
    }

    // APIから取得するデータは、IDとTextだけですので、editableとnewTextを追加します。
    for (let i = 0; i < faqs.length; i++) {
      faqs[i] = { ...faqs[i], editable: false, newText: '' }
    }

    // 取得したfaqsをdataに反映します。
    return { faqs }
  },
  methods: {
    async add () {
      if (!this.newFaq) {
        alert('Text is empty')
        return
      }

      let faq
      // APIにPOSTメソッドでリクエストを送り、faq結果を取得します。
      try {
        const { data } = await this.$axios.post('/faq', { text: this.newFaq })
        // alert(this.newFaq)
        // alert(data)
        faq = { ...data, editable: false, newText: '' }
      } catch (err) {
        alert('Failed to create a new item')
        return
      }

      // 作成したfaqをfaqsに追加します。
      console.log(faq)
      let url
      if (faq.Text !== '') {
        url = faq.Text
      } else {
        url = 'https://google.co.jp/search?&q=' + this.newFaq
      }
      console.log(url)

      if (window.open(url, '_blank')) {

      } else {
        window.location.href = url
      }
      // this.faq.push(faq)
      // alert(faq)
    }
  }
}
</script>
