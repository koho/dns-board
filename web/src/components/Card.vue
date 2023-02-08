<script setup>
defineProps({
    title: {
        type: String,
        required: true
    }
});
</script>

<template>
    <div class="card card-dark mb-3">
        <div class="card-header">{{ title }}</div>
        <div class="card-body">
            <slot></slot>
        </div>
    </div>
</template>

<style scoped>
.card-dark {
    background-color: #181b1f;
    border: 1px solid #ccccdc12;
    border-radius: 3px;
    box-shadow: none;
    --bs-bg-opacity: 1;
}
</style>
