// The message to show for a response that reports a failure.
//
// When the failure was on the server's side, the response carries a
// reference: the ID of the request, which finds the failure in the backend's
// logs. It is shown with the message so that it can be quoted in a report.
export function apiErrorMessage(
  result: { error?: string; request_id?: string } | null,
  status: number,
): string {
  const message = result?.error ?? `The backend returned an error (${status}).`;
  return result?.request_id
    ? `${message} (reference: ${result.request_id})`
    : message;
}
