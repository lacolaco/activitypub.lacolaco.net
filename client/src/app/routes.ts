import { Routes } from '@angular/router';
import { UserComponent } from './user/user.component';
import { SearchComponent } from './search/search.component';
import { requireAuthentication } from './core/auth/guard';

export const routes: Routes = [
  {
    path: 'users/:username',
    component: UserComponent,
  },
  {
    path: '@:username',
    redirectTo: 'users/:username',
  },
  {
    path: 'search',
    component: SearchComponent,
    canMatch: [requireAuthentication()],
  },
];
