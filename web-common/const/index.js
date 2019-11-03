let server = '';
console.log(process.env.NODE_ENV);
if (process.env.NODE_ENV == 'production') {
  server = 'http://172.16.0.3/';
} else {
  server = 'http://localhost:8888/';
}
const testRunInterval = 15;
const resultList = [
  {
    code: 0,
    name: '等待中',
    description: ''
  },
  {
    code: 1,
    name: '等待重判',
    description: ''
  },
  {
    code: 2,
    name: '编译中',
    description: ''
  },
  {
    code: 3,
    name: '运行并评判',
    description: ''
  },
  {
    code: 4,
    name: '正确',
    description: '答案正确，请再接再厉。'
  },
  {
    code: 5,
    name: '格式错误',
    description: '答案基本正确，但是格式不对。'
  },
  {
    code: 6,
    name: '答案错误',
    description:
      '答案不对，仅仅通过样例数据的测试并不一定是正确答案，一定还有你没想到的地方。'
  },
  {
    code: 7,
    name: '时间超限',
    description:
      '运行超出时间限制，检查下是否有死循环，或者应该有更快的计算方法。'
  },
  {
    code: 8,
    name: '内存超限',
    description: '超出内存限制，数据可能需要压缩，检查内存是否有泄露。'
  },
  {
    code: 9,
    name: '输出超限',
    description: '输出超过限制，你的输出比正确答案长了两倍，一定是哪里弄错了。'
  },
  {
    code: 10,
    name: '运行错误',
    description:
      '运行时错误，可能的原因有非法的内存访问，数组越界，指针漂移，调用禁用的系统函数。'
  },
  {
    code: 11,
    name: '编译错误',
    description: '编译错误，请点击后获得编译器的详细输出。'
  },
  {
    code: 12,
    name: '编译成功',
    description: ''
  },
  {
    code: 13,
    name: '测试运行',
    description: ''
  }
];
const langList = [
  'C',
  'Cpp',
  'Pascal',
  'Java',
  'Ruby',
  'Bash',
  'Python',
  'PHP',
  'Perl',
  'CSharp',
  'Objc',
  'FreeBasic',
  'Schema',
  'Clang',
  'Clang++',
  'Lua',
  'Javascript',
  'Go'
];

export { server, testRunInterval, resultList, langList };
