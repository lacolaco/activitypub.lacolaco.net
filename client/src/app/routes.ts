import { Routes, UrlSegment } from '@angular/router';
import { UserComponent } from './user/user.component';

export const routes: Routes = [
  {
    // "@:username" is not supported by default
    matcher: (url) => {
      const match = url[0]?.path.match(/^@([a-zA-Z0-9_]+)$/);
      if (match) {
        return {
          consumed: url,
          posParams: { username: new UrlSegment(match[1], {}) },
        };
      }
      return null;
    },
    component: UserComponent,
  },
  {
    path: '**',
    redirectTo: '@lacolaco',
  },
];
