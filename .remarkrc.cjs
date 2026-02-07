module.exports = {
  plugins: [
    './.opencode/node_modules/remark-preset-lint-recommended/index.js',
    './.opencode/node_modules/remark-lint-heading-increment/index.js',
    './.opencode/node_modules/remark-lint-no-duplicate-headings/index.js',
    './.opencode/node_modules/remark-lint-no-undefined-references/index.js',
    ['./.opencode/node_modules/remark-lint-list-item-indent/index.js', 'one'],
    ['./.opencode/node_modules/remark-lint-fenced-code-marker/index.js', '`'],
    [
      './.opencode/node_modules/remark-lint-no-missing-blank-lines/index.js',
      {exceptTightLists: true}
    ]
  ]
}
