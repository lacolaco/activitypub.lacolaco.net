import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Auth, idToken } from '@angular/fire/auth';
import { switchMap } from 'rxjs';
import { environment } from 'src/environments/environment';

export function authInterceptor(): HttpInterceptorFn {
  return (req, next) => {
    const auth = inject(Auth);
    // skip if not a request to our API
    if (!req.url.startsWith(environment.backend)) {
      return next(req);
    }

    return idToken(auth).pipe(
      switchMap((token) => {
        if (!token) {
          return next(
            req.clone({
              withCredentials: true,
            }),
          );
        }
        const authReq = req.clone({
          withCredentials: true,
          setHeaders: {
            Authorization: `Bearer ${token}`,
          },
        });
        return next(authReq);
      }),
    );
  };
}
