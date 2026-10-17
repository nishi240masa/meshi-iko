// ユーザがTimeSlotで選択できる時間帯の範囲を30分刻みで作成する関数
export function createTimeScale(startTime: string, endTime: string): string[] {
  const timeScale: string[] = [];
  const [startHour, startMinute] = startTime.split(":").map(Number);
  const [endHour, endMinute] = endTime.split(":").map(Number);

  let currentHour = startHour;
  let currentMinute = startMinute;

  while (currentHour < endHour || (currentHour === endHour && currentMinute <= endMinute)) {
    const formattedTime = `${String(currentHour).padStart(2, "0")}:${String(currentMinute).padStart(2, "0")}`;
    timeScale.push(formattedTime);
    currentMinute += 30;
    if (currentMinute >= 60) {
      currentMinute -= 60;
      currentHour++;
    }
  }

  return timeScale;
}
