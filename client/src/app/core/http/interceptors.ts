import { HttpInterceptorFn } from '@angular/common/http';
import { environment } from '../../../environments/environment';

export const withCredentialsInterceptor = (): HttpInterceptorFn => {
  return (req, next) => {
    // skip if not a request to our API
    const url = new URL(req.url);
    if (url.origin !== environment.backend) {
      return next(req);
    }
    return next(req.clone({ withCredentials: true }));
  };
};
