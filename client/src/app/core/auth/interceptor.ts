import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { Auth, user } from '@angular/fire/auth';
import { map, switchMap } from 'rxjs';

export function authInterceptor(): HttpInterceptorFn {
  return (req, next) => {
    const auth = inject(Auth);
    // skip if not a request to our API
    if (!req.url.startsWith('/api')) {
      return next(req);
    }

    return user(auth).pipe(
      map((user) => user?.getIdToken()),
      switchMap((token) => {
        if (!token) {
          return next(req);
        }
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
