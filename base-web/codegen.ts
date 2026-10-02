import type { CodegenConfig } from '@graphql-codegen/cli';

const config: CodegenConfig = {
  schema: '../base-engine/gen/schema.graphqls',
  documents: ['src/**/*.tsx', 'src/**/*.ts', 'src/**/*.graphql', '!src/__generated__/**'],
  generates: {
    './src/__generated__/': {
      preset: 'client',
      plugins: [],
      config: {
        useTypeImports: true,
      },
      presetConfig: {
        gqlTagName: 'gql',
      },
    },
  },
  ignoreNoDocuments: true,
};

export default config;
