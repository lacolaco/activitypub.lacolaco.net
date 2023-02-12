import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Auth } from '@angular/fire/auth';
import { from, switchMap } from 'rxjs';

export function authInterceptor(): HttpInterceptorFn {
  return (req, next) => {
    const auth = inject(Auth);
    // skip if not a request to our API
    if (!req.url.startsWith('/api')) {
      return next(req);
    }
    if (!auth.currentUser) {
      return next(req);
    }

    return from(auth.currentUser.getIdToken()).pipe(
      switchMap((token) => {
        const authReq = req.clone({
          setHeaders: {
            Authorization: `Bearer ${token}`,
          },
        });
        return next(authReq);
      }),
    );
  };
}
