import { Routes } from '@angular/router';
import { FooComponent } from './foo/foo.component';
import { SearchComponent } from './search/search.component';

export const routes: Routes = [
  {
    path: 'users/:username',
    component: FooComponent,
  },
  {
    path: 'search',
    component: SearchComponent,
  },
];
