import { Routes, UrlMatchResult, UrlSegment } from '@angular/router';
import { UserComponent } from './user/user.component';
import { SearchComponent } from './search/search.component';
import { requireAuthentication } from './core/auth/guard';

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
    path: 'search',
    component: SearchComponent,
    canMatch: [requireAuthentication()],
  },
  {
    path: '**',
    redirectTo: '@lacolaco',
  },
];
