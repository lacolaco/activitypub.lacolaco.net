/**
 * @see https://schema.org/PropertyValue
 */
export type PropertyValue = {
  type: 'PropertyValue';
  name: string;
  value: string;
};

export function buildPropertyValue(name: string, value: string): PropertyValue {
  return {
    type: 'PropertyValue',
    name,
    value,
  };
}
