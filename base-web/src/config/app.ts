const engineHttpUrl = 'http://localhost:1980';

/** Web 应用直接使用的后端连接配置。 */
export const appConfig = {
  engine: {
    httpUrl: engineHttpUrl,
    graphQLHttpUrl: `${engineHttpUrl}/graphql`,
    graphQLWebSocketUrl: `${engineHttpUrl.replace(/^http/, 'ws')}/graphql`,
    systemInitializationUrl: `${engineHttpUrl}/api/system-initialization`,
  },
} as const;
