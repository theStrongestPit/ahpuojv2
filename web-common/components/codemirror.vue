<template lang="pug">
    codemirror(v-model="mycode", :options="cmOptions", @input="onCmCodeChange")
</template>

<script>
// language
import 'codemirror/mode/clike/clike.js';
import 'codemirror/mode/go/go.js';
import 'codemirror/mode/python/python.js';
// theme css
import 'codemirror/theme/monokai.css';
//auto close brackets
import 'codemirror/addon/edit/closebrackets.js';
// require active-line.js
import 'codemirror/addon/selection/active-line.js';
// styleSelectedText
import 'codemirror/addon/selection/mark-selection.js';
import 'codemirror/addon/search/searchcursor.js';
// highlightSelectionMatches
import 'codemirror/addon/scroll/annotatescrollbar.js';
import 'codemirror/addon/search/matchesonscrollbar.js';
import 'codemirror/addon/search/searchcursor.js';
import 'codemirror/addon/search/match-highlighter.js';
// keyMap
import 'codemirror/mode/clike/clike.js';
import 'codemirror/addon/edit/matchbrackets.js';
import 'codemirror/addon/comment/comment.js';
import 'codemirror/addon/dialog/dialog.js';
import 'codemirror/addon/dialog/dialog.css';
import 'codemirror/addon/search/searchcursor.js';
import 'codemirror/addon/search/search.js';
import 'codemirror/keymap/sublime.js';
// foldGutter
import 'codemirror/addon/fold/foldgutter.css';
import 'codemirror/addon/fold/brace-fold.js';
import 'codemirror/addon/fold/comment-fold.js';
import 'codemirror/addon/fold/foldcode.js';
import 'codemirror/addon/fold/foldgutter.js';
import 'codemirror/addon/fold/indent-fold.js';
import 'codemirror/addon/fold/markdown-fold.js';
import 'codemirror/addon/fold/xml-fold.js';

export default {
  props: {
    code: {
      type: String,
      default: ''
    },
    language: {
      type: Number,
      default: 1
    },
    readonly: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      map: new Map([
        [0, 'text/x-csrc'],
        [1, 'text/x-c++src'],
        [3, 'text/x-java'],
        [6, 'text/x-python'],
        [17, 'text/x-go']
      ]),
      cmOptions: {
        tabSize: 4,
        styleActiveLine: false,
        lineNumbers: true,
        styleSelectedText: false,
        autoCloseBrackets: true,
        line: true,
        foldGutter: true,
        readOnly: false,
        gutters: ['CodeMirror-linenumbers', 'CodeMirror-foldgutter'],
        highlightSelectionMatches: {showToken: /\w/, annotateScrollbar: true},
        mode: 'text/x-csrc',
        keyMap: 'sublime',
        matchBrackets: true,
        showCursorWhenSelecting: true,
        theme: 'monokai',
        extraKeys: {Ctrl: 'autocomplete'}
      },
      mycode: this.code
    };
  },
  watch: {
    language() {
      console.log(this.map.get(this.language));
      this.cmOptions.mode = this.map.get(this.language);
      console.log(this.cmOptions.mode);
    }
  },
  mounted() {
    const self = this;
    self.cmOptions.readOnly = self.readonly;
    console.log('123');
    console.log(self.cmOptions);
    setTimeout(() => {
      (self.styleSelectedText = true), (self.styleActiveLine = true);
    }, 1800);
  },
  methods: {
    onCmCodeChange(newCode, event) {
      this.$emit('update:code', newCode);
    }
  }
};
</script>
<style lang="scss">
.CodeMirror-focused .cm-matchhighlight {
  background-image: url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAYAAABytg0kAAAAFklEQVQI12NgYGBgkKzc8x9CMDAwAAAmhwSbidEoSQAAAABJRU5ErkJggg==);
  background-position: bottom;
  background-repeat: repeat-x;
}
.cm-matchhighlight {
  background-color: rgba(lightblue, 0.4);
}
.CodeMirror-selection-highlight-scrollbar {
  background-color: green;
}
.CodeMirror {
  line-height: 20px;
  height: 550px;
}
.CodeMirror-sizer {
  margin-left: 50px !important;
}
</style>
