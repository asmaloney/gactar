// http://eslint.org/docs/user-guide/configuring

module.exports = {
  root: true,
  env: {
    es2021: true,
  },
  plugins: ['html', 'vue'],
  extends: ['plugin:vue/recommended', 'eslint:recommended', '@vue/prettier'],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
  },
}
