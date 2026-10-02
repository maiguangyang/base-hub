import { ApolloLink, Observable } from '@apollo/client/core';
import { getMainDefinition } from '@apollo/client/utilities';
import { publishActionFeedback } from '@/lib/actionFeedback';
import { actionErrorReason } from './actionErrors';
import { actionLabel } from './actionMessages';

/** Every write receives feedback; composite saves explicitly publish their final outcome. */
export function createActionFeedbackLink(): ApolloLink {
  return new ApolloLink((operation, forward) => {
    const definition = getMainDefinition(operation.query);
    if (definition.kind !== 'OperationDefinition' || definition.operation !== 'mutation' || operation.getContext().adminFeedback === false) return forward(operation);
    const fields = definition.selectionSet.selections.filter((item) => item.kind === 'Field');
    const label = actionLabel(fields[0]?.name.value ?? '');
    return new Observable((observer) => {
      let result: ApolloLink.Result | undefined;
      let reportedError = false;
      const subscription = forward(operation).subscribe({
        next: (next) => {
          result = next;
          if (next.errors?.length && !reportedError) {
            reportedError = true;
            publishActionFeedback('error', `${label}失败：${actionErrorReason(next.errors[0])}`);
          }
          observer.next(next);
        },
        error: (cause: unknown) => {
          if (!reportedError) publishActionFeedback('error', `${label}失败：${actionErrorReason(cause)}`);
          observer.error(cause);
        },
        complete: () => {
          if (reportedError) { observer.complete(); return; }
          if (fields.length && fields.every((field) => result?.data?.[field.alias?.value ?? field.name.value] != null && result.data[field.alias?.value ?? field.name.value] !== false)) {
            publishActionFeedback('success', `${label}成功。`);
          } else publishActionFeedback('error', `${label}失败：操作未完成，服务端未返回有效结果，请核对后重试。`);
          observer.complete();
        },
      });
      return () => subscription.unsubscribe();
    });
  });
}
