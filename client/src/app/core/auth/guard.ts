import { inject } from '@angular/core';
import { Auth, user } from '@angular/fire/auth';
import { CanMatchFn, Router } from '@angular/router';
import { map } from 'rxjs';

export function requireAuthentication(): CanMatchFn {
  return (route, state) => {
    const auth = inject(Auth);
    const router = inject(Router);

    return user(auth).pipe(
      map((user) => {
        if (!user) {
          return router.createUrlTree(['/']);
        }
        return true;
      }),
    );
  };
}
