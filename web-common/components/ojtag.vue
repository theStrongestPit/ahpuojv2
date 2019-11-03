<template lang="pug">
    span(class="oj-tag",:class="[typeClass]")
      slot
</template>

<script>
export default {
  name: 'OjTag',
  props: {
    type: {
      type: String,
      validator(val) {
        return ['primary', 'success', 'danger', 'warnning'].indexOf(val);
      },
      default: ''
    }
  },
  computed: {
    typeClass() {
      return this.type ? `oj-tag--${this.type}` : '';
    }
  }
};
</script>