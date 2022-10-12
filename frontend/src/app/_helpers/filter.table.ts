export function advancedFilterFn(data: any, searchTerm: string): boolean {
  let [phrase, columns] = searchTerm.split(' in:').map((str) => str.trim());
  return Object.keys(data).some((key: any) => {
    if (columns?.length) {
      let result;
      columns.split(',').forEach((column) => {
        if (
          column.toLowerCase().trim() === key.toLowerCase() &&
          data[key].toLowerCase().includes(phrase.toLowerCase())
        ) {
          result = true;
        }
      });
      return result;
    }
    if (data[key] && !columns?.length) {
      return JSON.stringify(data).toLowerCase().includes(phrase.toLowerCase());
    }
    return [];
  });
}
