const maxGraphqlInt = 2_147_483_647n;

export function yuanToFen(input: string): number {
  const value = input.trim();
  if (!/^(0|[1-9]\d*)(\.\d{1,2})?$/.test(value)) {
    throw new Error('INVALID_YUAN_AMOUNT');
  }
  const [yuan, fraction = ''] = value.split('.');
  const fen = BigInt(yuan) * 100n + BigInt(fraction.padEnd(2, '0'));
  if (fen > maxGraphqlInt) {
    throw new Error('GRAPHQL_INT_OVERFLOW');
  }
  return Number(fen);
}

export function fenToYuan(input: number): string {
  if (!Number.isInteger(input) || input < 0 || input > Number(maxGraphqlInt)) {
    throw new Error('INVALID_FEN_AMOUNT');
  }
  const value = BigInt(input);
  return (value / 100n).toString() + '.' + String(value % 100n).padStart(2, '0');
}
