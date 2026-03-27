export default async function Page({
  params,
}: {
  params: Promise<{ param: string }>;
}) {
  const { param } = await params;

  return (
    <div className="flex items-center justify-center text-xl bg-amber-50 text-black h-screen">
      {param}
      hiiiiiii
    </div>
  );
}
