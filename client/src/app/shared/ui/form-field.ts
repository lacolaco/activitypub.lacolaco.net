import { ChangeDetectionStrategy, Component, ContentChild, Directive, Input, NgModule } from '@angular/core';

let nextUniqueId = 0;

@Directive({
  selector: 'input:not([type=checkbox]), textarea, select',
  standalone: true,
  host: {
    class: `
    block p-2 w-full text-base text-gray-900 bg-white rounded border border-solid border-gray-400 appearance-none placeholder:text-transparent peer
    focus:outline-none focus:ring-0 focus:border-purple-700 
    disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed`,
    '[attr.id]': 'id',
  },
})
export class FormControlDirective {
  @Input('aria-labelledby') id = `app-form-control-${nextUniqueId++}`;
}

@Component({
  selector: 'app-form-field',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <ng-content></ng-content>
    <label
      [style.display]="label && showLabel ? undefined : 'none'"
      [attr.aria-hidden]="'false'"
      [attr.for]="inputId"
      class="absolute text-xs text-gray-500 peer-focus:text-purple-700 origin-[0] bg-white px-1 top-[16px] -translate-y-4 left-1"
    >
      {{ label }}
    </label>
  `,
  host: {
    '[class.show-label]': 'showLabel',
  },
  styles: [
    `
      :host {
        display: block;
        position: relative;
      }
      :host-context(.show-label) {
        padding-top: 0.5em;
      }
    `,
  ],
})
export class FormFieldComponent {
  @Input() showLabel = false;
  @Input() label = '';
  @ContentChild(FormControlDirective) formControl?: FormControlDirective;

  get inputId() {
    return this.formControl?.id;
  }
}

@NgModule({
  imports: [FormFieldComponent, FormControlDirective],
  exports: [FormFieldComponent, FormControlDirective],
})
export class FormFieldModule {}
