<template>
  <div class="container">
    <div class="main">
      <UCard class="box-card">
        <template #header>
          <div class="card-header">
            <h1>Welcome to kk_scheduler</h1>
          </div>
        </template>
        <p>
          This is a frontend interface for kk_scheduler.
        </p>
        <p>
          You can manage your tasks and schedules here.
        </p>
        <div style="margin-top: 20px;">
          <UButton 
            variant="link" 
            color="primary" 
            to="https://github.com/cruvie/kk_scheduler"
            target="_blank"
          >
            Visit GitHub Repository
          </UButton>
        </div>
      </UCard>
    </div>
  </div>
</template>

<style scoped>
.container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.main {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.box-card {
  width: 480px;
  text-align: center;
}

.card-header {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>