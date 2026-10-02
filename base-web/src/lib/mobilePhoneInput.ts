/** 与账号服务当前接受的 11 位大陆手机号规则保持一致。 */
export const mobilePhoneInputProps = {
  inputMode: 'numeric' as const,
  autoComplete: 'tel',
  maxLength: 11,
  minLength: 11,
  pattern: '1[34578][0-9]{9}',
  title: '请输入11位有效手机号',
};
