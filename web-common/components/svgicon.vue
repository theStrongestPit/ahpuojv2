<template lang="pug">
  svg(:class="iconClass",aria-hidden="true")
    use(:xlink:href="iconName")
</template>

<script>
export default {
  name: 'SvgIcon',
  props: {
    name: {
      type: String,
      required: true
    },
    className: {
      type: String,
      default: ''
    }
  },
  computed: {
    iconName() {
      return `#icon-${this.name}`;
    },
    iconClass() {
      if (this.className) {
        return 'icon ' + this.className;
      } else {
        return 'icon';
      }
    }
  }
};
</script>

<style lang="scss" scoped>
.icon {
  display: inline-block;
  width: 40px;
  height: 40px;
  vertical-align: middle;
  fill: currentColor;
  overflow: hidden;
}
</style>