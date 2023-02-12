import { Directive } from '@angular/core';

@Directive({
  standalone: true,
  host: {
    class: 'disabled:bg-gray-200 disabled:cursor-not-allowed disabled:text-gray-400',
  },
})
export class AppDisabledButton {}

@Directive({
  selector: 'button[app-icon-button]',
  standalone: true,
  host: {
    class: 'flex justify-center text-center rounded-full h-12 w-12 p-3 hover:bg-gray-100',
  },
  hostDirectives: [AppDisabledButton],
})
export class AppIconButton {}

@Directive({
  selector: 'button[app-stroked-button]',
  standalone: true,
  host: {
    class:
      'flex justify-center text-center rounded py-2 px-4 border border-solid border-gray-300 leading-none hover:bg-gray-100',
  },
  hostDirectives: [AppDisabledButton],
})
export class AppStrokedButton {}
